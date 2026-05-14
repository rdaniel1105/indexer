package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirectoryReaderError(t *testing.T) {
	c := require.New(t)

	path := "$5asdas"
	expectedErr := fmt.Errorf("open %v: no such file or directory", path)

	files, err := DirectoryReader(path)
	c.EqualError(err, expectedErr.Error())
	c.Nil(files)
}

func TestDirectoryReader(t *testing.T) {
	c := require.New(t)

	files, err := DirectoryReader("testdata/maildir")
	c.NoError(err)
	c.NotNil(files)
	c.NotEmpty(files)
}
