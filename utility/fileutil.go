package utility

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
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
	err := os.MkdirAll(path, fs.ModeDir)
	if err != nil {
		log.Panicf("Couldn't create folder %s!\n", path)
	}
}
