package fs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

type FileUtils struct {
}

func (f FileUtils) Read(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func (f FileUtils) HomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}
