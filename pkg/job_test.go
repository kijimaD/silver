package silver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobRun(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	task1 := NewTask(
		"Run uname command1",
		buf,
	)
	{
		instfunc := func() error {
			err := task1.Exec("uname")

			return err
		}
		task1.instCmd = instfunc
	}
	task2 := NewTask(
		"Run uname command2",
		buf,
	)
	{
		instfunc := func() error {
			err := task2.Exec("uname")

			return err
		}
		task2.instCmd = instfunc
	}

	job := NewJob([]Task{task1, task2})
	job.Run()

	expect := `[Run uname command1]
  => [exec] uname
  => Linux
=> Success install
[Run uname command2]
  => [exec] uname
  => Linux
=> Success install
`
	assert.Equal(t, expect, buf.String())
}
