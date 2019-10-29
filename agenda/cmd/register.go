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
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a user",
	Long: `Register a user and set its name, password and email,  For example:
		agenda register -u UserTest -p your_password -e your_email   or 
		agenda register --user=UserTest --password=your_password --email=your_email`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		emailaddr, _ := cmd.Flags().GetString("email")
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
		if len(emailaddr) == 0 {
			fmt.Println("Error: Email address must be set")
			cmd.Help()
			return
		}
		//validation of email address
		matched, _ := regexp.MatchString(`[\w-]+@[\w]+(?:\.[\w]+)+`, emailaddr)
		if matched == false {
			fmt.Println("Error: Your email address is invalid, please check")
			return
		}
		//check if username was unique
		if fileObjR, errR := os.OpenFile("users.txt", os.O_RDONLY|os.O_CREATE, 0644); errR == nil {
			defer fileObjR.Close()
			if contents, err := ioutil.ReadAll(fileObjR); err == nil {
				result := strings.Replace(string(contents), "\n", "", 0)
				infos := strings.Split(result, "\n")
				for i := 0; i < len(infos); i += 3 {
					if username == infos[i] {
						fmt.Println("Error: This username has been used, please choose another one")
						return
					}
				}
			}
		}
		//write to log
		fileObjW, err := os.OpenFile("users.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Failed to open the file", err.Error())
			os.Exit(2)
		}
		defer fileObjW.Close()
		content := username + "\n" + password + "\n" + emailaddr + "\n"
		if _, err := io.WriteString(fileObjW, content); err == nil {
			fmt.Println(username + " registered successfully with email: " + emailaddr)
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	registerCmd.Flags().StringP("user", "u", "", "Help message for username")
	registerCmd.Flags().StringP("password", "p", "", "Help message for password")
	registerCmd.Flags().StringP("email", "e", "", "Help message for email")
}
