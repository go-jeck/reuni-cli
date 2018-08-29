package helper

import (
	"fmt"
	"log"
	"os"
)

type CredentialCaller interface {
	Write()
}

type CredentialHelper struct {
	Token        string
	RefreshToken string
}

func (f *CredentialHelper) Write() {
	os.Mkdir("/var/tmp/reuni", os.ModePerm)
	f.writeFile(f.Token, "/var/tmp/reuni/token")
	f.writeFile(f.RefreshToken, "/var/tmp/reuni/refresh")
	fmt.Println("Credentials created")
}

func (f *CredentialHelper) RewriteToken() {
	f.writeFile(f.Token, "/var/tmp/reuni/token")
	fmt.Println("Credentials refreshed")
}

func (f *CredentialHelper) writeFile(content, location string) {
	keyfile, err := os.Create(location)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer keyfile.Close()
	fmt.Fprintf(keyfile, "%v", content)
}
