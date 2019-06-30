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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "parsec-ec2",
	Short: "Start and stop Parsec EC2 instances with a single command",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var installPath, region, cfgFile, goPath, instanceType, projectPath, amiName string
var volumeSize int

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "aws region")
	RootCmd.PersistentFlags().StringVarP(&instanceType, "instance-type", "i", "", "ec2 instance type")
	RootCmd.PersistentFlags().StringVarP(&amiName, "ami-name", "a", "", "AMI Name to create ec2 instance")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		goPath = os.Getenv("GOPATH")

		projectPath = fmt.Sprintf("%s/src/github.com/tereshkin/parsec-ec2", goPath)
		installPath = fmt.Sprintf("%s/.parsec-ec2", home)

		// Search config in home directory with name ".parsec-ec2" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".parsec-ec2")
	}

	viper.SetEnvPrefix("parsec_ec2")
	viper.AutomaticEnv() // read in environment variables that match
	serverKey = viper.GetString("server_key")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
