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

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialisation command to be run after initial installs and subsequent upgrades",
	Long: `
Copies the latest Terraform and Windows provisioning templates to the
$HOME/.parsec-ec2 directory and runs 'terraform init' to initialise any
plugins required by Terraform.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the install directory exists
		fmt.Println("Checking for existing installation...")
		if _, err := os.Stat(installPath); os.IsNotExist(err) {
			fmt.Print("No existing installation found. Copying templates and initialising... ")
			// If it doesn't exist, make it
			err := os.Mkdir(installPath, 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Print("Existing installation found. Copying latest templates... ")
		}

		tSrc := fmt.Sprintf("%s/%s", projectPath, Template)
		tDst := fmt.Sprintf("%s/%s", installPath, Template)
		if err := cp(tSrc, tDst); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		uSrc := fmt.Sprintf("%s/%s", projectPath, Userdata)
		uDst := fmt.Sprintf("%s/%s", installPath, Userdata)
		if err := cp(uSrc, uDst); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		init := tfCmd([]string{TfCmdInit})

		if err := execute(init); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Complete.")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
