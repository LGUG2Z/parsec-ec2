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

	"encoding/json"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running Parsec EC2 instance",
	Long: `
Stops a Parsec EC2 instance created using the start command. Under the
hood this command runs 'terraform destroy', with removes all AWS resources
that are identified for creation in the terraform template.

This command depends on session information that is created by the start
command and stored in $HOME/.parsec-ec2/currentSession.json, so if this
has been manually modified or removed after running the start command,
the stop command will not execute. In this situation it is still possible
to manually run 'terraform destroy'. You will receive prompts for variable
values, but these can all be left blank with the exception of the region
variable, which can be set to the region the instances were started in.

Example:

parsec-ec2 stop
`,
	Run: func(cmd *cobra.Command, args []string) {
		session := fmt.Sprintf("%s/%s", installPath, CurrentSession)

		bytes, err := ioutil.ReadFile(session)
		if err != nil {
			fmt.Println("No session information found.")
			os.Exit(1)
		}

		var p TfVars

		if err := json.Unmarshal(bytes, &p); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		destroy := tfCmdVars(p, []string{TfCmdDestroy, TfFlagForce, TfFlagApprove})

		fmt.Println("Terminating all AWS resources created by this session... ")
		if err := execute(destroy); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := os.Remove(session); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("All resources have been successfully terminated.")
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
