package helpers

import (
	"io/ioutil"
	"os"
)

// DirectoryReader reads the content inside the given path and return them in a string array.
func DirectoryReader(root string) ([]string, error) {
	files := make([]string, 0)

	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return files, nil
}

// DirectoryChecker checks if the path given is a directory.
func DirectoryChecker(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
