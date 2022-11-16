package helpers

import (
	"fmt"
	"testing"
)

func TestCreateEmailStructError(t *testing.T) {
	path := "$5asdas"

	expectedErr := fmt.Errorf("file reading error: open %v: no such file or directory", path)

	_, _, err := CreateEmailStruct(path)

	if err == nil {
		t.Error(expectedErr.Error())
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error FAILED: expected %v got %v\n", expectedErr, err)
	}

}

func TestCreateEmailStruct(t *testing.T) {
	path := "../../enron_mail_20110402/maildir/allen-p/sent/1"

	_, _, err := CreateEmailStruct(path)

	if err != nil {
		t.Error("FAILED: expecting err to be nil")
	}
}
