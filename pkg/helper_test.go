package syscheck

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExistCmd(t *testing.T) {
	assert.True(t, isExistCmd("ls"))
	assert.False(t, isExistCmd("THIS_IS_INVALID_COMMAND"))
}

func TestIsExistFile(t *testing.T) {
	assert.True(t, isExistFile("."))
	assert.False(t, isExistFile("THIS_IS_INVALID_FILE"))
}

func TestIsExistFileHomeDir(t *testing.T) {
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir
	expect := fmt.Sprintf("%s/Project", homeDir)
	result, err := expandTilde("~/Project")
	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}
