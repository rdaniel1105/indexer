package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEmailStructError(t *testing.T) {
	c := require.New(t)

	path := "$5asdas"
	expectedErr := fmt.Errorf("file reading error: open %v: no such file or directory", path)

	_, _, err := CreateEmailStruct(path)
	c.EqualError(err, expectedErr.Error())
}

func TestCreateEmailStruct(t *testing.T) {
	c := require.New(t)

	email, repeated, err := CreateEmailStruct("testdata/maildir/sample/sent/1")
	c.NoError(err)
	c.False(repeated)
	c.NotNil(email)
	c.Equal("alice@example.com", email.From)
	c.Equal("bob@example.com", email.To)
	c.Equal("Sample message", email.Subject)
}
