package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func Edit(config map[string]string) map[string]string {
	configJSON, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err.Error())
	}

	tempfileHelper := TempFileHelper{
		Payload: string(configJSON),
	}
	tempfileHelper.Write()

	cmd := exec.Command("vi", tempfileHelper.Filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	result := make(map[string]string)
	err = json.Unmarshal(tempfileHelper.Open(), &result)
	if err != nil {
		fmt.Println(err.Error())
	}

	tempfileHelper.Close()

	return result
}
