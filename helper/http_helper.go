package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HttpCaller interface {
	SendRequest() (*http.Response, error)
}

type HttpHelper struct {
	URL           string
	Method        string
	Authorization string
	RefreshToken  string
	Payload       []byte
}

func (h *HttpHelper) SendRequest() (*http.Response, error) {
	client := &http.Client{}
	body := bytes.NewBuffer(h.Payload)
	req, err := http.NewRequest(h.Method, h.URL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+h.Authorization)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		req.Header.Set("Authorization", "Bearer "+h.generateNewToken())
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func FetchData(caller HttpCaller, data interface{}) error {
	resp, err := caller.SendRequest()
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("HTTP Error: " + resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(data)
	resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func (h *HttpHelper) generateNewToken() string {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/generateToken", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+h.RefreshToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	res := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&res)
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	credentialHelper := CredentialHelper{
		Token: fmt.Sprint(res["token"]),
	}
	credentialHelper.RewriteToken()

	return res["token"]
}
