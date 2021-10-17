package fs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

type ReadFn func(string) ([]byte, error)

func Read(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func HomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}
