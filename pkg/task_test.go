package silver

import "testing"

func TestTaskRun(t *testing.T) {
	testfunc := func() error {
		err := Run("uname")
		return err
	}

	task := NewTask(
		"Run uname command",
		[]errorFunc{testfunc, testfunc},
	)
	task.Run()
}
