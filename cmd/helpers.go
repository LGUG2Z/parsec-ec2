package cmd

import (
	"fmt"
	"net/http"
	"os"

	"io/ioutil"
	"os/exec"

	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

func getExternalIP() (string, error) {
	resp, err := http.Get("https://ipv4.icanhazip.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", nil
		}
		return fmt.Sprintf("%s/%s", strings.TrimSpace(string(b)), "32"), nil
	}

	return "", errors.New("Could not get external ip address.")
}

func getEc2Client(region string) (*ec2.EC2, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	config := aws.Config{
		Region: aws.String(region),
	}

	return ec2.New(session, &config), nil
}

func getVpcID(svc *ec2.EC2) (string, error) {
	vpc, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{})

	if err != nil {
		return "", err
	}

	if len(vpc.Vpcs) < 1 {
		fmt.Println(`You have at some point manually deleted the default VPC created by AWS for this region.
parsec-ec2 will not function for this region until a default VPC is recreated for it.
You can still try to launch Parsec EC2 instances in other regions that have retained
their default VPC.`)
		os.Exit(0)
	}

	return *vpc.Vpcs[0].VpcId, nil
}

func getSubnetID(svc *ec2.EC2, availabilityZone string) (string, error) {
	values := []*string{&availabilityZone}

	filter := ec2.Filter{
		Name:   aws.String("availability-zone"),
		Values: values,
	}

	filters := []*ec2.Filter{&filter}

	describeSubnetsInput := ec2.DescribeSubnetsInput{
		Filters: filters,
	}

	result, err := svc.DescribeSubnets(&describeSubnetsInput)

	if err != nil {
		return "", err
	}

	if len(result.Subnets) == 0 {
		fmt.Printf("Could not get the subnet id for availability zone %s.\n", availabilityZone)
		os.Exit(1)
	}

	return *result.Subnets[0].SubnetId, nil
}

func copy(source, destination string) error {
	b, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destination, b, 0644)
}

func hasServerKey(serverKey string) bool {
	return len(serverKey) > 0
}

func tfCmd(args []string) *exec.Cmd {
	command := exec.Command(Terraform, args...)
	command.Dir = installPath
	command.Env = os.Environ()

	return command
}

func tfCmdVars(p TfVars, args []string) *exec.Cmd {
	command := exec.Command(Terraform, args...)

	command.Dir = installPath

	command.Env = os.Environ()
	command.Env = append(command.Env, fmt.Sprintf("TF_VAR_instance_type=%s", p.InstanceType))
	command.Env = append(command.Env, fmt.Sprintf("TF_VAR_region=%s", p.Region))
	command.Env = append(command.Env, fmt.Sprintf("TF_VAR_server_key=%s", p.ServerKey))
	command.Env = append(command.Env, fmt.Sprintf("TF_VAR_spot_price=%s", p.SpotPrice))
	command.Env = append(command.Env, fmt.Sprintf("TF_VAR_subnet_id=%s", p.SubnetID))
	command.Env = append(command.Env, fmt.Sprintf("TF_VAR_vpc_id=%s", p.VpcID))
	command.Env = append(command.Env, fmt.Sprintf("TF_VAR_ami=%s", p.AmiID))
	command.Env = append(command.Env, fmt.Sprintf("TF_VAR_volume_size=%d", p.VolumeSize))
	command.Env = append(command.Env, fmt.Sprintf("TF_VAR_ip=%s", p.IP))

	return command
}

func executeSilent(command *exec.Cmd) error {
	commandErr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	err = command.Start()
	if err != nil {
		return err
	}

	errOutput, err := ioutil.ReadAll(commandErr)
	if err != nil {
		return err
	}

	if len(errOutput) > 0 {
		return fmt.Errorf("Error executing Terraform command: %s\nError Output: %s", command.Args[1], errOutput)
	}

	return nil
}

func executeReturn(command *exec.Cmd) ([]byte, error) {
	initOut, err := command.StdoutPipe()
	if err != nil {
		return []byte{}, err
	}

	initErr, err := command.StderrPipe()
	if err != nil {
		return []byte{}, err
	}

	err = command.Start()
	if err != nil {
		return []byte{}, err
	}

	stdOutput, err := ioutil.ReadAll(initOut)
	if err != nil {
		return []byte{}, err
	}

	errOutput, err := ioutil.ReadAll(initErr)
	if err != nil {
		return []byte{}, err
	}

	if len(errOutput) > 0 {
		return []byte{}, fmt.Errorf("Error executing Terraform command: %s\nError Output: %s", command.Args[1], errOutput)
	}

	return stdOutput, nil
}

func executePrint(command *exec.Cmd) error {
	commandOut, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	commandErr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	err = command.Start()
	if err != nil {
		return err
	}

	stdOutput, err := ioutil.ReadAll(commandOut)
	if err != nil {
		return err
	}

	errOutput, err := ioutil.ReadAll(commandErr)
	if err != nil {
		return err
	}

	if len(errOutput) > 0 {
		return fmt.Errorf("Error executing Terraform command: %s\nError output: %s", command.Args[1], errOutput)
	}

	fmt.Printf("%s\n", stdOutput)
	fmt.Printf("%s\n", errOutput)

	return nil
}
