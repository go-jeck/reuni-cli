package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-squads/reuni-cli/helper"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

var serviceName string

type service struct {
	Name               string     `json:"name"`
	AuthorizationToken string     `json:"authorization_token"`
	CreatedAt          *time.Time `json:"created_at"`
	OrganizationId     int        `json:"organization_id"`
	CreatedBy          string     `json:"created_by"`
}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage service in your organization",
	Long:  `Manage service in your organization. Organization name is required input for this operation`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !strings.EqualFold(cmd.CalledAs(), "service") {
			key = getToken()
			if len(organizationName) < 1 {
				fmt.Println("Invalid Organization")
				os.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listServiceCmd = &cobra.Command{
	Use:   "list",
	Short: "Display all service",
	Run: func(cmd *cobra.Command, args []string) {
		httphelper := &helper.HttpHelper{
			URL:           "http://127.0.0.1:8080/" + organizationName + "/services",
			Method:        "GET",
			Authorization: key,
		}

		var services []service
		err := helper.FetchData(httphelper, &services)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		output := []string{
			"#|Service Name|Authorization Token|Created At|Created By",
		}
		for i, s := range services {
			output = append(output, fmt.Sprintf("%v|%v|%v|%v|%v", i+1, s.Name, s.AuthorizationToken, s.CreatedAt, s.CreatedBy))
		}
		result := columnize.SimpleFormat(output)
		fmt.Println(result)
	},
}

var createServiceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new service",
	Long:  `This command is for creating new service. use flag -o for organization name and -s for service name`,
	Run: func(cmd *cobra.Command, args []string) {
		data := make(map[string]string)
		data["name"] = serviceName
		dataJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		httphelper := &helper.HttpHelper{
			URL:           "http://127.0.0.1:8080/" + organizationName + "/services",
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

var deleteServiceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete existing service",
	Long:  `This command is for deleting a service. use flag -o for organization name and -s for service name`,
	Run: func(cmd *cobra.Command, args []string) {
		data := make(map[string]string)
		data["name"] = serviceName
		dataJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		httphelper := &helper.HttpHelper{
			URL:           "http://127.0.0.1:8080/" + organizationName + "/services",
			Method:        "DELETE",
			Authorization: key,
			Payload:       dataJSON,
		}

		res, err := httphelper.SendRequest()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if res.StatusCode == http.StatusOK {
			fmt.Println("Organization Deleted")
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
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(listServiceCmd)
	serviceCmd.AddCommand(createServiceCmd)
	serviceCmd.AddCommand(deleteServiceCmd)

	serviceCmd.PersistentFlags().StringVarP(&organizationName, "organization", "o", "", "Your organization name")
	serviceCmd.PersistentFlags().StringVarP(&serviceName, "service", "s", "", "Your service name")
}
