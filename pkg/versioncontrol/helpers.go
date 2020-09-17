package versioncontrol

import (
	"log"
	"os"
)

func clearWorkingDir(path string) error {

	err := os.RemoveAll(path)
	if err != nil {
		log.Fatal(err)
	}

	return os.MkdirAll(path, 0755)
}
