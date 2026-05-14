package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirectoryCheckerError(t *testing.T) {
	c := require.New(t)

	path := "$5asdas"
	expectedErr := fmt.Errorf("stat %v: no such file or directory", path)

	isDir, err := DirectoryChecker(path)
	c.EqualError(err, expectedErr.Error())
	c.False(isDir)
}

func TestDirectoryCheckDir(t *testing.T) {
	c := require.New(t)

	isDir, err := DirectoryChecker("testdata/maildir")
	c.NoError(err)
	c.True(isDir)
}

func TestDirectoryCheckFile(t *testing.T) {
	c := require.New(t)

	isDir, err := DirectoryChecker("testdata/maildir/sample/sent/1")
	c.NoError(err)
	c.False(isDir)
}
