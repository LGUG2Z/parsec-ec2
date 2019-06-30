package cmd

import (
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Terraform command
const Terraform = "terraform"

// Terraform CLI Commands
const (
	TfCmdApply   = "apply"
	TfCmdDestroy = "destroy"
	TfCmdInit    = "init"
	TfCmdOutput  = "output"
	TfCmdPlan    = "plan"
	TfCmdRefresh = "refresh"
)

// Parsec Terraform Template Outputs
const (
	TfOutputInstanceType   = "instance_type"
	TfOutputRegion         = "region"
	TfOutputServerKey      = "server_key"
	TfOutputSpotInstanceID = "spot_instance_id"
	TfOutputSpotPrice      = "spot_price"
	TfOutputSubnetID       = "subnet_id"
	TfOutputVpcID          = "vpc_id"
)

// Terraform CLI Command Flags
const (
	TfFlagForce = "-force"
	TfFlagJSON  = "-json"
	TfFlagAutoApprove = "-auto-approve"
)

// Filenames
const (
	Template       = "parsec.tf"
	Userdata       = "user_data.tmpl"
	CurrentSession = "currentSession.json"
)

// Product Description and Instance Statuses
const (
	Windows = "Windows"
	OK      = "ok"
)

func ec2Regions() map[string]endpoints.Region {
	partition := endpoints.AwsPartition()
	services := partition.Services()
	ec2 := services[endpoints.Ec2ServiceID]
	return ec2.Regions()
}

func isValidRegion(validRegions map[string]endpoints.Region, input string) bool {
	for _, valid := range validRegions {
		if input == valid.ID() {
			return true
		}
	}
	return false
}

func gInstances() []string {
	return []string{
		ec2.InstanceTypeG22xlarge,
		ec2.InstanceTypeG28xlarge,
		ec2.InstanceTypeG34xlarge,
		ec2.InstanceTypeG3sXlarge,
		ec2.InstanceTypeG38xlarge,
		ec2.InstanceTypeG316xlarge,
	}
}

func isValidGInstance(validInstances []string, input string) bool {
	for _, valid := range validInstances {
		if input == valid {
			return true
		}
	}
	return false
}

func gMinVolumeSize() int {
	return 35
}

func isValidVolumeSize(minSize, input int) bool {
	return input >= minSize
}
