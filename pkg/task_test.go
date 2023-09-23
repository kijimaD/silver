package silver

import (
	"bytes"
	"os"
	"testing"
)

func TestTaskRun(t *testing.T) {
	testfunc := func() error {
		err := Run("uname", os.Stdout)
		return err
	}

	buf := &bytes.Buffer{}
	task := NewTask(
		"Run uname command",
		[]boolFunc{},
		[]errorFunc{testfunc, testfunc},
		buf,
	)
	task.Run()
}

func TestTaskNotMet(t *testing.T) {
	depsfunc := func() bool {
		return IsExistCmd("not_found_cmd")
	}

	testfunc := func() error {
		err := Run("uname", os.Stdout)
		return err
	}

	buf := &bytes.Buffer{}
	task := NewTask(
		"Run uname command",
		[]boolFunc{depsfunc},
		[]errorFunc{testfunc, testfunc},
		buf,
	)
	task.Run()
	// TODO: bufをチェックする
}
