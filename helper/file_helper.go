package helper

import (
	"fmt"
	"log"
	"os"
)

type FileCaller interface {
	writeFile()
}

type FileHelper struct {
	Payload string
}

func (f *FileHelper) WriteFile() {
	os.Mkdir("/var/tmp/reuni", os.ModePerm)
	keyfile, err := os.Create("/var/tmp/reuni/key")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer keyfile.Close()
	fmt.Fprintf(keyfile, "%v", f.Payload)
	fmt.Println("Credentials created")
}
