/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "user login",
	Long: `User login agenda using username and password. For example:
		agenda login -u your_username -p your_password    or
		agenda login --user=your_username --password=your_password`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		if len(username) == 0 {
			fmt.Println("Error: Username must be set")
			cmd.Help()
			return
		}
		if len(password) == 0 {
			fmt.Println("Error: Password must be set")
			cmd.Help()
			return
		}
		//check if password for username was correct
		if fileObjR, errR := os.OpenFile("users.txt", os.O_RDONLY|os.O_CREATE, 0644); errR == nil {
			defer fileObjR.Close()
			if contents, err := ioutil.ReadAll(fileObjR); err == nil {
				result := strings.Replace(string(contents), "\n", "", 0)
				infos := strings.Split(result, "\n")
				for i := 0; i < len(infos); i += 3 {
					if username == infos[i] {
						if password == infos[i+1] {
							fmt.Println("Succeed:  user " + username + " login successfully")
							return
						}
						fmt.Println("Fail: password incorrect, please try again")
						return
					}
				}
				fmt.Println("Fail: you are not registered, please register first")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().StringP("user", "u", "", "Help message for username")
	loginCmd.Flags().StringP("password", "p", "", "Help message for password")
}
