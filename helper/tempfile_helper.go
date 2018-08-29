package helper

import (
	"io/ioutil"
	"log"
	"os"
)

type TempFileCaller interface {
	Write()
	Open()
	Close()
}

type TempFileHelper struct {
	Payload  string
	Filename string
	file     *os.File
}

func (f *TempFileHelper) Write() {
	tempFile, err := ioutil.TempFile("", "reuni")
	f.file = tempFile
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.file.Write([]byte(f.Payload)); err != nil {
		log.Fatal(err)
	}

	f.Filename = f.file.Name()
}

func (f *TempFileHelper) Open() []byte {
	bytes, err := ioutil.ReadFile(f.Filename)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func (f *TempFileHelper) Close() {
	if err := f.file.Close(); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Filename)
}
