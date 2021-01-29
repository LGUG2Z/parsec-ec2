// Copyright © 2017 Jade Iqbal <jadeiqbal@fastmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new Parsec EC2 instance",
	Long: `
Makes a spot request for the requested EC2 instance type in the specified region.
If PARSEC_EC2_SERVER_KEY has not been exported in the shell rc file, it must be
passed to the command using the --server-key flag.

The amount to bid relative to the current highest spot price for an instance is
specified using the --bid flag, so if the current highest spot price is $0.20,
running the command with --bid 0.10 will send a spot request with a bid price
of $0.30.

If the --plan flag is used, the spot request will not be sent and instead the
'terraform plan' command will be run which will output to the console the details
of any AWS resources that will be created by running the start command.

Examples:

parsec-ec2 start --aws-region eu-west-1 --instance-type g3.4xlarge --bid 0.10
parsec-ec2 start --aws-region eu-west-1 --instance-type g2.2xlarge --bid 0.10 ---server-key xxxxx
parsec-ec2 start --aws-region eu-central-1 --instance-type g2.2xlarge --bid 0.10 --plan
`,
	Run: func(cmd *cobra.Command, args []string) {
		if !hasServerKey(serverKey) {
			fmt.Println("Add 'export PARSEC_EC2_SERVER_KEY=xxxxx' to your shell rc file or use the --server-key flag.")
			os.Exit(1)
		}

		if !isValidRegion(ec2Regions(), region) {
			fmt.Printf("%s is not a valid AWS region id.\n", region)
			os.Exit(1)
		}

		if !isValidGInstance(gInstances(), instanceType) {
			fmt.Printf("%s is not a valid EC2 GPU instance type id.\n", instanceType)
			os.Exit(1)
		}

		if !isValidVolumeSize(gMinVolumeSize(), volumeSize) {
			fmt.Printf("%d is less than a minumum volume size %d.\n", volumeSize, gMinVolumeSize())
			os.Exit(1)
		}

		ec2Client, err := getEc2Client(region)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var p TfVars

		if err := p.Calculate(ec2Client, region, serverKey, instanceType, amiName, volumeSize); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var start *exec.Cmd

		// TODO: Use a template to generate a .tfvars file
		if plan {
			start = tfCmdVars(p, []string{TfCmdPlan})

			fmt.Printf(
				"Planning spot request for a %s "+
				"instance from AMI %s in %s with volume size %d "+
				"and a bid of $%s...\n\n",
				p.InstanceType, p.AmiName, p.Region, p.VolumeSize, p.SpotPrice)
			
			if err := executePrint(start); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("If you are happy with this plan run the start command again without the --plan flag.")
		} else {
			start = tfCmdVars(p, []string{TfCmdApply, TfFlagAutoApprove})
			fmt.Printf(
				"Making spot request for a %s "+
				"instance from AMI %s in %s with volume size %d "+
				"and a bid of $%s...\n\n",
				p.InstanceType, p.AmiName, p.Region, p.VolumeSize, p.SpotPrice)

			if err := executePrint(start); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Spot request made successfully. Check the status of the spot request with 'parsec-ec2 status'.")

			if err := p.Write(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func calculateUserBid(cheapestSpotPrice string, bidIncrease float64) string {
	spotPrice, _ := strconv.ParseFloat(cheapestSpotPrice, 64)
	userBid := spotPrice + bidIncrease

	return fmt.Sprint(userBid)
}

var (
	bid       float64
	serverKey string
	plan      bool
)

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().Float64VarP(&bid, "bid", "b", 0.00, "amount to bid relative to the current highest spot price")
	startCmd.Flags().StringVarP(&serverKey, "server-key", "k", "", "Parsec server key")
	startCmd.Flags().BoolVarP(&plan, "plan", "p", false, "plan out the resources to be created without creating them")
	startCmd.Flags().IntVarP(&volumeSize, "volume-size", "s", 0, "volume size to be used for EC2 instance")
}
