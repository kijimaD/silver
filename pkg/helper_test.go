package silver

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExistCmd(t *testing.T) {
	assert.True(t, IsExistCmd("ls"))
	assert.False(t, IsExistCmd("THIS_IS_INVALID_COMMAND"))
}

func TestIsExistFile(t *testing.T) {
	assert.True(t, IsExistFile("."))
	assert.False(t, IsExistFile("THIS_IS_INVALID_FILE"))
}

func TestIsExistFileHomeDir(t *testing.T) {
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir
	expect := fmt.Sprintf("%s/Project", homeDir)
	result, err := expandTilde("~/Project")

	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}
