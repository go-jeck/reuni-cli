package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-squads/reuni-cli/helper"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

type configView struct {
	Version       int               `json:"version"`
	Configuration map[string]string `json:"configuration"`
	Created_by    string            `json:"created_by"`
}

type versionView struct {
	Version int `json:"version"`
}

var version int
var configKey, configVal string

var configurationCmd = &cobra.Command{
	Use:   "configuration",
	Short: "Manage Configuration of your namespace",
	Long:  `Manage Configuration of your namespace. Organization, service and namespace name are required`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !strings.EqualFold(cmd.CalledAs(), "namespace") {
			key = getToken()
			refreshToken = getRefreshToken()
			if len(organizationName) < 1 {
				fmt.Println("Invalid Organization")
				os.Exit(1)
			}

			if len(serviceName) < 1 {
				fmt.Println("Invalid Service")
				os.Exit(1)
			}

			if len(namespaceName) < 1 {
				fmt.Println("Invalid Namespace")
				os.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listConfigurationCmd = &cobra.Command{
	Use:   "list",
	Short: "Display version list of configuration",
	Run: func(cmd *cobra.Command, args []string) {
		displayListVersions()
	},
}

var displayConfugurationCmd = &cobra.Command{
	Use:   "show",
	Short: "Display detail of configuration",
	Run: func(cmd *cobra.Command, args []string) {
		displayConfig()
	},
}

var updateAllConfigurationCmd = &cobra.Command{
	Use:   "update",
	Short: "Update configuration",
	Run: func(cmd *cobra.Command, args []string) {
		updateAllConfig()
	},
}

var setConfigurationCmd = &cobra.Command{
	Use:   "set",
	Short: "Update configuration",
	Run: func(cmd *cobra.Command, args []string) {
		setConfig()
	},
}

var unsetConfigurationCmd = &cobra.Command{
	Use:   "unset",
	Short: "Update configuration",
	Run: func(cmd *cobra.Command, args []string) {
		unsetConfig()
	},
}

var changesConfigurationCmd = &cobra.Command{
	Use:   "changes",
	Short: "Display changes from configuration",
	Run: func(cmd *cobra.Command, args []string) {
		displayChanges()
	},
}

func init() {
	rootCmd.AddCommand(configurationCmd)
	configurationCmd.AddCommand(listConfigurationCmd)
	configurationCmd.AddCommand(displayConfugurationCmd)
	configurationCmd.AddCommand(updateAllConfigurationCmd)
	configurationCmd.AddCommand(setConfigurationCmd)
	configurationCmd.AddCommand(unsetConfigurationCmd)
	configurationCmd.AddCommand(changesConfigurationCmd)

	configurationCmd.PersistentFlags().StringVarP(&organizationName, "organization", "o", "", "Your organization name")
	configurationCmd.PersistentFlags().StringVarP(&serviceName, "service", "s", "", "Your service name")
	configurationCmd.PersistentFlags().StringVarP(&namespaceName, "namespace", "n", "", "Your namespace name")
	configurationCmd.PersistentFlags().StringVarP(&configurationsData, "configurations", "c", "", "Your configurations")
	configurationCmd.PersistentFlags().IntVarP(&version, "versions", "v", 0, "Version")
	configurationCmd.PersistentFlags().StringVarP(&configKey, "key", "", "", "Configuration key")
	configurationCmd.PersistentFlags().StringVarP(&configVal, "value", "", "", "Configuration value")
}

//base func
func displayConfig() {
	if version < 1 {
		version = fetchLatestVersion(organizationName, serviceName, namespaceName)
	}

	config := fetchConfiguration(organizationName, serviceName, namespaceName, version)

	fmt.Println(displayHeader(serviceName, namespaceName, config.Version))
	fmt.Println(displayKeyVal(config.Configuration))
}

func unsetConfig() {
	version = fetchLatestVersion(organizationName, serviceName, namespaceName)
	config := fetchConfiguration(organizationName, serviceName, namespaceName, version)

	newConfig := config.Configuration
	if _, exists := newConfig[configKey]; !exists {
		fmt.Println("key not found!")
		return
	}

	delete(newConfig, configKey)

	data := make(map[string]interface{})
	data["configuration"] = newConfig
	dataJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	updateConfig(dataJSON)
}

func setConfig() {
	version = fetchLatestVersion(organizationName, serviceName, namespaceName)
	config := fetchConfiguration(organizationName, serviceName, namespaceName, version)

	newConfig := config.Configuration

	newConfig[configKey] = configVal

	data := make(map[string]interface{})
	data["configuration"] = newConfig
	dataJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	updateConfig(dataJSON)
}

func displayListVersions() {
	httphelper := &helper.HttpHelper{
		URL:           fmt.Sprintf("%v/%v/%v/%v/versions", "http://127.0.0.1:8080", organizationName, serviceName, namespaceName),
		Method:        "GET",
		Authorization: key,
		RefreshToken:  refreshToken,
	}

	data := make(map[string][]int)
	err := helper.FetchData(httphelper, &data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	output := []string{
		"Service|" + serviceName,
		"Namespace|" + namespaceName,
		"Versions|" + strings.Trim(strings.Replace(fmt.Sprint(data["versions"]), " ", ",", -1), "[]"),
	}
	result := columnize.SimpleFormat(output)
	fmt.Println(result)
}

func displayChanges() {
	if version < 1 {
		version = fetchLatestVersion(organizationName, serviceName, namespaceName)
	} else {
		if version == 1 {
			displayConfig()
			return
		}
	}

	previousVersion := version - 1

	currentConfig := fetchConfiguration(organizationName, serviceName, namespaceName, version)
	previousConfig := fetchConfiguration(organizationName, serviceName, namespaceName, previousVersion)

	curConfig := currentConfig.Configuration
	prevConfig := previousConfig.Configuration

	deletedConfig := getConfigurationDeleted(curConfig, prevConfig)
	createdConfig := getConfigurationCreated(curConfig, prevConfig)
	changedConfig := getConfigurationChanged(curConfig, prevConfig)

	templateHeader := []string{
		"Service|" + serviceName,
		"Namespace|" + namespaceName,
		"Version|" + strconv.Itoa(version),
		"Changed By|" + currentConfig.Created_by,
	}

	fmt.Println(columnize.SimpleFormat(templateHeader))

	if len(deletedConfig) > 0 {
		fmt.Println("Deleted Keys")
		fmt.Println(displayKeyVal(deletedConfig))
	}
	if len(createdConfig) > 0 {
		fmt.Println("Created Keys")
		fmt.Println(displayKeyVal(createdConfig))
	}
	if len(changedConfig) > 0 {
		fmt.Println("Changed Keys")
		fmt.Println(displayKeyVal(changedConfig))
	}
}

func updateAllConfig() {
	version = fetchLatestVersion(organizationName, serviceName, namespaceName)
	config := fetchConfiguration(organizationName, serviceName, namespaceName, version)

	newConfig := helper.Edit(config.Configuration)

	data := make(map[string]interface{})
	data["configuration"] = newConfig
	dataJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	updateConfig(dataJSON)
}

func getConfigurationDeleted(currentConfig, previousConfig map[string]string) map[string]string {
	return getKeydifference(previousConfig, currentConfig)
}
func getConfigurationCreated(currentConfig, previousConfig map[string]string) map[string]string {
	return getKeydifference(currentConfig, previousConfig)
}

func getConfigurationChanged(currentConfig, previousConfig map[string]string) map[string]string {
	mapB := map[string]bool{}
	for k, _ := range previousConfig {
		mapB[k] = true
	}

	res := make(map[string]string)
	for k, v := range currentConfig {
		if _, ok := mapB[k]; ok {
			if previousConfig[k] != v {
				res[k] = "from '" + previousConfig[k] + "' to '" + v + "'"
			}
		}
	}

	return res
}

func getKeydifference(configA, configB map[string]string) map[string]string {
	mapB := map[string]bool{}
	for k, _ := range configB {
		mapB[k] = true
	}

	res := make(map[string]string)
	for k, v := range configA {
		if _, ok := mapB[k]; !ok {
			res[k] = v
		}
	}

	return res
}

//ui helper
func displayKeyVal(config map[string]string) string {
	templateBody := []string{
		"#|Key|Value",
	}

	i := 1
	for k, v := range config {
		templateBody = append(templateBody, fmt.Sprintf("%v|%v|%v", i, k, v))
		i++
	}

	return columnize.SimpleFormat(templateBody)
}

func displayHeader(serviceName, namespaceName string, version int) string {
	templateHeader := []string{
		"Service|" + serviceName,
		"Namespace|" + namespaceName,
		"Version|" + strconv.Itoa(version),
	}

	return columnize.SimpleFormat(templateHeader)
}

//model helper
func fetchLatestVersion(organizationName, serviceName, namespaceName string) int {
	httphelper := &helper.HttpHelper{
		URL:           fmt.Sprintf("%v/%v/%v/%v/latest", "http://127.0.0.1:8080", organizationName, serviceName, namespaceName),
		Method:        "GET",
		Authorization: key,
		RefreshToken:  refreshToken,
	}

	data := make(map[string]int)
	err := helper.FetchData(httphelper, &data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return data["version"]
}

func fetchConfiguration(organizationName, serviceName, namespaceName string, version int) configView {
	httphelper := &helper.HttpHelper{
		URL:           fmt.Sprintf("%v/%v/%v/%v/%v", "http://127.0.0.1:8080", organizationName, serviceName, namespaceName, version),
		Method:        "GET",
		Authorization: key,
		RefreshToken:  refreshToken,
	}

	var config configView
	err := helper.FetchData(httphelper, &config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return config
}

func updateConfig(payload []byte) {
	httphelper := &helper.HttpHelper{
		URL:           fmt.Sprintf("%v/%v/%v/%v", "http://127.0.0.1:8080", organizationName, serviceName, namespaceName),
		Method:        "POST",
		Authorization: key,
		RefreshToken:  refreshToken,
		Payload:       payload,
	}

	res, err := httphelper.SendRequest()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if res.StatusCode == http.StatusCreated {
		fmt.Println("New Configuration Created")
	} else {
		data := make(map[string]interface{})
		err = json.NewDecoder(res.Body).Decode(&data)
		res.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("HTTP Error " + fmt.Sprint(data["status"]) + ": " + fmt.Sprint(data["message"]))
	}
}
