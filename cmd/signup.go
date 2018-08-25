// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/go-squads/reuni-cli/helper"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// signupCmd represents the signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Creating new user",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var name, username, password, confirmPassword, email string
		fmt.Println("signup called")
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Name: ")
			name, _ = reader.ReadString('\n')
			name = string(name[0 : len(name)-1])
			if name != "" {
				break
			}
		}
		for {
			fmt.Print("Username: ")
			username, _ = reader.ReadString('\n')
			username = string(username[0 : len(username)-1])
			if username != "" {
				break
			}
		}
		for {
			fmt.Print("Password: ")
			bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
			password = string(bytePassword)
			fmt.Println()
			if password != "" {
				break
			}
		}
		for {
			fmt.Print("Confirm Password: ")
			byteConfirmPassword, _ := terminal.ReadPassword(int(syscall.Stdin))
			confirmPassword = string(byteConfirmPassword)
			fmt.Println()
			if confirmPassword != "" && strings.Compare(password, confirmPassword) == 0 {
				break
			}
		}
		for {
			fmt.Print("Email: ")
			email, _ = reader.ReadString('\n')
			email = string(email[0 : len(email)-1])
			if email != "" {
				break
			}
		}

		data := make(map[string]string)
		data["name"] = name
		data["username"] = username
		data["password"] = password
		data["email"] = email
		dataJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		httphelper := &helper.HttpHelper{
			URL:           "http://127.0.0.1:8080/signup",
			Method:        "POST",
			Authorization: key,
			Payload:       dataJSON,
		}

		res, err := httphelper.SendRequest()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if res.StatusCode == http.StatusCreated {
			fmt.Println("User Created")
		} else {
			data := make(map[string]interface{})
			err = json.NewDecoder(res.Body).Decode(&data)
			res.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("HTTP Error " + fmt.Sprint(data["status"]) + ": " + fmt.Sprint(data["message"]))
		}
	},
}

func init() {
	rootCmd.AddCommand(signupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
