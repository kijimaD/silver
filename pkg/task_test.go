package silver

import "testing"

func TestTaskRun(t *testing.T) {
	datefunc := func() error {
		err := Run("date")
		return err
	}

	task := NewTask(
		"this task running date command",
		[]func() error{datefunc, datefunc},
	)
	task.Run()
}
