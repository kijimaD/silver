package silver

import "testing"

func TestTaskRun(t *testing.T) {
	datefunc := func() {
		Run("date")
	}

	task := NewTask(
		"this task running date command",
		[]func(){datefunc, datefunc},
	)
	task.Run()
}
