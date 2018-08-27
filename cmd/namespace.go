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

var configurationsData string

type namespaceView struct {
	OrganizationId int       `json:"organization_id"`
	ServiceName    string    `json:"service_name"`
	Namespace      string    `json:"namespace"`
	ActiveVersion  int       `json:"version"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type namespaceStore struct {
	Namespace     string            `json:"namespace"`
	Configuration map[string]string `json:"configurations"`
}

var namespaceName string

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage namespace inside service",
	Long:  `Manage namespace inside service. Organization name and service name is required.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !strings.EqualFold(cmd.CalledAs(), "namespace") {
			key = getToken()
			if len(organizationName) < 1 {
				fmt.Println("Invalid Organization")
				os.Exit(1)
			}

			if len(serviceName) < 1 {
				fmt.Println("Invalid Service")
				os.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createNamespaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new namespace",
	Long: `This command is for creating new namespace. Use flag -o for
organization name, -s for service name and -n for namespace name. 
this command also require intial configuration using flag -c followed by
json object.

for example
reuni-cli namespace create -o org -s service -n default -c '{"firstKey":"firstVal","secondKey":"secondVal"}'
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(namespaceName) < 1 {
			fmt.Println("Invalid namespace name")
			return
		}

		if len(configurationsData) < 1 {
			fmt.Println("Configuration can't be empty")
			return
		}
		fmt.Println(configurationsData)

		configurationsPayload := make(map[string]string)
		err := json.Unmarshal([]byte(configurationsData), &configurationsPayload)
		fmt.Println(configurationsPayload)

		var data namespaceStore
		data.Namespace = namespaceName
		data.Configuration = configurationsPayload
		dataJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(dataJSON))

		httphelper := &helper.HttpHelper{
			URL:           "http://127.0.0.1:8080/" + organizationName + "/" + serviceName + "/namespaces",
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
			fmt.Println("Namespace Created")
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

var listNamespaceCmd = &cobra.Command{
	Use:   "list",
	Short: "Display all namespace",
	Run: func(cmd *cobra.Command, args []string) {
		httphelper := &helper.HttpHelper{
			URL:           "http://127.0.0.1:8080/" + organizationName + "/" + serviceName + "/namespaces",
			Method:        "GET",
			Authorization: key,
		}

		var namespaces []namespaceView
		err := helper.FetchData(httphelper, &namespaces)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		output := []string{
			"#|Namespace Name|Active Version|Created By|Created At|Updated At",
		}
		for i, n := range namespaces {
			output = append(output, fmt.Sprintf("%v|%v|%v|%v|%v|%v", i+1, n.Namespace, n.ActiveVersion, n.CreatedBy, n.CreatedAt, n.UpdatedAt))
		}
		result := columnize.SimpleFormat(output)
		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
	namespaceCmd.AddCommand(createNamespaceCmd)
	namespaceCmd.AddCommand(listNamespaceCmd)

	namespaceCmd.PersistentFlags().StringVarP(&organizationName, "organization", "o", "", "Your organization name")
	namespaceCmd.PersistentFlags().StringVarP(&serviceName, "service", "s", "", "Your service name")
	namespaceCmd.PersistentFlags().StringVarP(&namespaceName, "namespace", "n", "", "Your namespace name")
	namespaceCmd.PersistentFlags().StringVarP(&configurationsData, "configurations", "c", "", "Your configurations")
}
