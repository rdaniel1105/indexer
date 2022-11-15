package helpers

import (
	"io/ioutil"
)

// DirectoryReader reads the content inside the given path and return them in a string array.
func DirectoryReader(path string) ([]string, error) {
	files := make([]string, 0)

	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return files, nil
}
