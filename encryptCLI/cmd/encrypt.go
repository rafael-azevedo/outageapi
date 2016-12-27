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
	"errors"
	"fmt"

	"github.com/rafael-azevedo/HPOMOutageTool/utils"
	"github.com/spf13/cobra"
)

var (
	key     string
	keyFlag string
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypts your password for use with the outage tool ",
	Long: `If no enryption key is provided this command encypts the provided plaintext password and returns the key and encrypted password. 

If an encryption key is provided and it meets the criteria required it will be used to encyprt the password. The program will then return the encrypted password.

You can use an key supplied by encryptCLI generate.

Use this password in your app.toml for the OutageTool and store the key in an enviromental variable

	 ORACLE 	store in variable  	ORACLEDBKEY
	 RethinkDB  store in variable 	RETHINKDBKEY


example app.toml:

[oracle]
username = "user1"
#store your encrypted passwords here for oracle
password = "3ce758035a5806857c0ac66b42493782d85952fcaa75bfc8b0b88bd1"
hostname = "host.staples.com"
port = 	"301939"
servicename = "Database service name "

[rethinkdb]
username = ""
#store your encrypted password here for rethinkdb
password = "3ce758035a5806857c0ac66b42493782d85952fcaa75bfc8b0b88bd1"
server  = "127.0.0.1"
port	= "28015"
database = "outagetool"

[API]
server = "127.0.0.1"
port = "443"
sslcert = "HPOMOutageTool/API/outage.pem"
sslkey = "HPOMOutageTool/API/outage.key"

	`,

	RunE: encyptCMD,
}

func init() {
	RootCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	encryptCmd.Flags().StringVarP(&keyFlag, "Key", "k", "", "base64 encoded key")
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func encyptCMD(cmd *cobra.Command, args []string) error {
	var e utils.EncryptOutput
	if len(args) != 1 {
		return errors.New("password needs to be provided")
	}
	password := string(args[0])
	if keyFlag != "" {
		err := e.EncryptWithKey64(keyFlag, password)
		if err != nil {
			return err
		}
		fmt.Printf("Key: %s\nPassword: %s\n", e.Key, e.Password)
		return nil
	}
	err := e.GenKeyAndEncrypt(password)
	if err != nil {
		return err
	}

	fmt.Printf("Key: %s\nPassword: %s\n", e.Key, e.Password)
	return nil

}
