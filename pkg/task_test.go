package silver

import "testing"

func TestTaskRun(t *testing.T) {
	testfunc := func() error {
		err := Run("uname")
		return err
	}

	task := NewTask(
		"Run uname command",
		[]boolFunc{},
		[]errorFunc{testfunc, testfunc},
	)
	task.Run()
}

func TestTaskNotMet(t *testing.T) {
	depsfunc := func() bool {
		return IsExistCmd("not_found_cmd")
	}

	testfunc := func() error {
		err := Run("uname")
		return err
	}

	task := NewTask(
		"Run uname command",
		[]boolFunc{depsfunc},
		[]errorFunc{testfunc, testfunc},
	)
	task.Run()
}
