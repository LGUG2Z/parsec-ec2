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

// Terraform CLI Command Flags
const (
	TfFlagForce = "-force"
	TfFlagJSON  = "-json"
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
	return services[endpoints.Ec2ServiceID].Regions()
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
