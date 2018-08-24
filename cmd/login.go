package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/go-squads/reuni-cli/helper"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var username string
var password string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(username) <= 0 {
			fmt.Println("Invalid credentials provided")
			for {
				fmt.Print("Username: ")
				reader := bufio.NewReader(os.Stdin)
				username, _ = reader.ReadString('\n')
				username = string(username[0 : len(username)-1])
				if username != "" {
					break
				}
			}
		}

		fmt.Println("Your password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		password = string(bytePassword)

		data := make(map[string]string)
		data["username"] = username
		data["password"] = password
		dataJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		httphelper := &helper.HttpHelper{
			URL:           "http://127.0.0.1:8080/login",
			Method:        "POST",
			Authorization: "",
			Payload:       dataJSON,
		}
		res := make(map[string]interface{})
		err = helper.FetchData(httphelper, &res)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fileHelper := &helper.FileHelper{
			Payload: fmt.Sprint(res["token"]),
		}
		fileHelper.WriteFile()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Your username for login")
}
