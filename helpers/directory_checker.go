package helpers

import "os"

// DirectoryChecker checks if the path given is a directory.
func DirectoryChecker(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
