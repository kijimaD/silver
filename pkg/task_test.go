package silver

import (
	"bytes"
	"testing"
)

func TestTaskRun(t *testing.T) {
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
	// TODO: bufをチェックする
}
