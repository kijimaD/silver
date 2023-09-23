package silver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskInst(t *testing.T) {
	buf := &bytes.Buffer{}
	task := NewTask(
		"Run uname command",
		buf,
	)
	testfunc := func() error {
		err := task.Exec("uname")
		return err
	}

	task.instCmds = []errorFunc{testfunc}

	task.Run()

	expect := `[Run uname command]
  => [exec] uname
  => Linux
=> Success install
`
	assert.Equal(t, expect, buf.String())
}

func TestTaskNotMet(t *testing.T) {
	buf := &bytes.Buffer{}
	task := NewTask(
		"Run uname command",
		buf,
	)
	depsfunc := func() bool {
		return IsExistCmd("not_found_cmd")
	}
	testfunc := func() error {
		err := task.Exec("uname")
		return err
	}

	task.depsCmds = []boolFunc{depsfunc}
	task.instCmds = []errorFunc{testfunc, testfunc}

	task.Run()

	expect := `[Run uname command]
=> Dependencies not met, skip
`
	assert.Equal(t, expect, buf.String())
}

func TestTaskAlreadyAchived(t *testing.T) {
	buf := &bytes.Buffer{}
	task := NewTask(
		"Run uname command",
		buf,
	)
	targetfunc := func() bool {
		return IsExistCmd("uname")
	}
	testfunc := func() error {
		err := task.Exec("uname")
		return err
	}

	task.targetCmds = []boolFunc{targetfunc}
	task.instCmds = []errorFunc{testfunc, testfunc}

	task.Run()

	expect := `[Run uname command]
=> Already achieved, skip
`
	assert.Equal(t, expect, buf.String())
}
