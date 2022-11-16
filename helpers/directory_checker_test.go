package helpers

import (
	"fmt"
	"testing"
)

func TestDirectoryCheckerError(t *testing.T) {
	path := "$5asdas"

	expectedErr := fmt.Errorf("stat %v: no such file or directory", path)

	IsDir, err := DirectoryChecker(path)

	if err == nil {
		t.Error(expectedErr.Error())
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error FAILED: expected %v got %v\n", expectedErr, err)
	}

	if IsDir {
		t.Error("FAILED: Expecting IsDir to be false")
	}
}

func TestDirectoryCheckDir(t *testing.T) {
	path := "../../enron_mail_20110402"
	IsDir, err := DirectoryChecker(path)

	if err != nil {
		t.Error("FAILED: expecting err to be nil")
	}

	if !IsDir {
		t.Error("FAILED: expecting IsDir to be true")
	}
}

func TestDirectoryCheckFile(t *testing.T) {
	path := "../../enron_mail_20110402/maildir/allen-p/sent/1"
	IsDir, err := DirectoryChecker(path)

	if err != nil {
		t.Error("FAILED: expecting err to be nil")
	}

	if IsDir {
		t.Error("FAILED: expecting IsDir to be false")
	}
}
