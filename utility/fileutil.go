package utility

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Copy copies a file to another location
func Copy(src, dst string) (int64, error) {

	srcFile, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !srcFile.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	copiedBytes, err := io.Copy(destination, source)
	return copiedBytes, err
}

// CreateFolder Create a folder and check if creation was successful
func CreateFolder(path string) {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		log.Panicf("Couldn't create folder %s!", path)
	} else {
		log.Infof("Successfully created %s!", path)
	}
}
