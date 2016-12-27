// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

	"github.com/rafael-azevedo/HPOMOutageTool/utils"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a private key in base64 for your outage tool",
	Long: `This option generates a key in base64 that can be used to encyrpt passwords to store in enviroment or decrypt stored password from base64 encoding.

	 For example:

./encryptCLI generate 
wfqUaWNtw3Y6HA2m2UclBA==

You can then store the key "UtRBLPlh5Nkaih/OgUj1Gg==" in an enviromental variable for the OutageTool.

ORACLEBDKEY=UtRBLPlh5Nkaih/OgUj1Gg==

This ORACLEDBKEY will be used to decrypt the password for ORACLE in app.toml on the Outage Tool.

suppose your password is "somepassword"

you have generated and encrypted the password to 3ce758035a5806857c0ac66b42493782d85952fcaa75bfc8b0b88bd1 with key UtRBLPlh5Nkaih/OgUj1Gg==

example app.toml:

[oracle]
username = "user1"
password = "3ce758035a5806857c0ac66b42493782d85952fcaa75bfc8b0b88bd1"
hostname = "host.staples.com"
port = 	"301939"
servicename = "Database service name "

[rethinkdb]
username = ""
password = ""
server  = "127.0.0.1"
port	= "28015"
database = "outagetool"

[API]
server = "127.0.0.1"
port = "443"
sslcert = "HPOMOutageTool/API/outage.pem"
sslkey = "HPOMOutageTool/API/outage.key"


when you run this through decryption with 
	key UtRBLPlh5Nkaih/OgUj1Gg== 
	password 3ce758035a5806857c0ac66b42493782d85952fcaa75bfc8b0b88bd1 

the password will be decrypted to "somepassword"

`,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := utils.GenString(16)
		if err != nil {
			fmt.Printf("There was an error generating the key %s\n", err)
		}
		fmt.Println(key)
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
