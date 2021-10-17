package fs

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

type ReadFn func(string) ([]byte, error)
type WriteFn func(string, []byte, fs.FileMode) error

func Read(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func Write(file string, data []byte, perm fs.FileMode) error {
	return ioutil.WriteFile(file, data, perm)
}

func HomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}
