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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the initialisation status of a launched EC2 instance",
	Long: `
Queries the launched instance and gets the current initialisation status.

Once an instance is reporting a status of initialised, it may still take
some time for the instance to show up in the Parsec desktop application.
This is because time is still required for the provisioning script to run
on the instance, which is what will allow the Parsec application to launch
and log in with the provided Parsec server key.
`,
	Run: func(cmd *cobra.Command, args []string) {
		session := fmt.Sprintf("%s/%s", installPath, CurrentSession)

		bytes, err := ioutil.ReadFile(session)
		if err != nil {
			fmt.Println("There are no sessions currently running.")
			os.Exit(0)
		}

		var p TfVars

		if err := json.Unmarshal(bytes, &p); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ec2Client, err := getEc2Client(p.Region)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		refresh := tfCmdVars(p, []string{TfCmdRefresh})
		if err := execute(refresh); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var o TfOutputs
		if err := o.Read(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(o.SpotInstanceID.Value) < 1 {
			fmt.Println("The spot instance request is awaiting fulfilment.")
			os.Exit(0)
		}

		instanceIds := []*string{&o.SpotInstanceID.Value}

		describeInstanceStatusInput := ec2.DescribeInstanceStatusInput{
			InstanceIds: instanceIds,
		}

		describeInstanceStatusOutput, err := ec2Client.DescribeInstanceStatus(&describeInstanceStatusInput)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(describeInstanceStatusOutput.InstanceStatuses) < 1 {
			if o.SpotBidStatus.Value == "instance-terminated-by-price" {
				fmt.Println("The spot price rose above your bid price and your instance was terminated. Run 'parsec-ec2 stop' to cleanup.")
				os.Exit(0)
			}
			fmt.Println("The spot instance request has been filled but the instance initialisation status is not available yet.")
			os.Exit(0)
		}

		instanceStatus := describeInstanceStatusOutput.InstanceStatuses[0].InstanceStatus.Status

		if *instanceStatus == OK {
			fmt.Println("The instance has been initialised.")
			fmt.Println("It will be connectable once the provisioning script has finished running.")
		} else {
			fmt.Println("The instance is initialising.")
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
