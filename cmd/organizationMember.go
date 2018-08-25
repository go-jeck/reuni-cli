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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-squads/reuni-cli/helper"
	"github.com/spf13/cobra"
)

// organizationMemberCmd represents the organizationMember command
var organizationMemberCmd = &cobra.Command{
	Use:   "organizationMember",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("organizationMember called")
		cmd.Help()
	},
}

// organizationMemberCmd represents the organizationMember command
var addOrganizationMemberCmd = &cobra.Command{
	Use:   "add",
	Short: "add member to organization",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		data := make(map[string]string)
		data["name"] = organizationName
		dataJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		httphelper := &helper.HttpHelper{
			URL:           "http://127.0.0.1:8080/organization",
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
			fmt.Println("Organization Created")
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
	rootCmd.AddCommand(organizationMemberCmd)
	organizationMemberCmd.AddCommand(addOrganizationMemberCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// organizationMemberCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// organizationMemberCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
