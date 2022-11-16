package helpers

import (
	"fmt"
	"testing"
)

func TestDirectoryReaderError(t *testing.T) {
	path := "$5asdas"

	expectedErr := fmt.Errorf("open %v: no such file or directory", path)

	files, err := DirectoryReader(path)

	if err == nil {
		t.Error(expectedErr.Error())
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error FAILED: expected %v got %v\n", expectedErr, err)
	}

	if files != nil {
		t.Error("FAILED: Expecting files to be nil")
	}
}

func TestDirectoryReader(t *testing.T) {
	path := "../../enron_mail_20110402"
	files, err := DirectoryReader(path)

	if err != nil {
		t.Error("FAILED: expecting err to be nil")
	}

	if files == nil {
		t.Error("FAILED: expecting files array to exist")
	}
}
