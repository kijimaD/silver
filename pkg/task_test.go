package silver

import "testing"

func TestTaskRun(t *testing.T) {
	testfunc := func() error {
		err := Run("uname")
		return err
	}

	task := NewTask(
		"this task running date command",
		[]errorFunc{testfunc, testfunc},
	)
	task.Run()
}
