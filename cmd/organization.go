package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-squads/reuni-cli/helper"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

type Organization struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var organizationName string
var key string

var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "Manage your organization",
	Long:  `Manage your organization`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		key = getToken()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("organization called")
	},
}

var listOrganizationCmd = &cobra.Command{
	Use:   "list",
	Short: "Creating new Organization",
	Long:  `Creating new Organization. to create new organization, use flag -n for organization name`,
	Run: func(cmd *cobra.Command, args []string) {
		httphelper := &helper.HttpHelper{
			URL:           "http://127.0.0.1:8080/organization",
			Method:        "GET",
			Authorization: key,
		}

		var organizations []Organization
		err := helper.FetchData(httphelper, &organizations)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		output := []string{
			"#|Organization Name|Your Role",
		}
		for i, o := range organizations {
			output = append(output, fmt.Sprintf("%v|%v|%v", i+1, o.Name, o.Role))
		}
		result := columnize.SimpleFormat(output)
		fmt.Println(result)
	},
}

var createOrganizationCmd = &cobra.Command{
	Use:   "create",
	Short: "Creating new Organization",
	Long:  `Creating new Organization. to create new organization, use flag -n for organization name`,
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
	rootCmd.AddCommand(organizationCmd)
	organizationCmd.AddCommand(listOrganizationCmd)
	organizationCmd.AddCommand(createOrganizationCmd)

	createOrganizationCmd.Flags().StringVarP(&organizationName, "Organization name", "o", "", "Name for your new organization")
}
