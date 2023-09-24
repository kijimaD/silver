package silver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobRun(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	task1 := NewTask("Run uname command1", WithWriter(buf))
	task1.SetFuncs(ExecFuncParam{
		TargetCmd: nil,
		DepCmd:    nil,
		InstCmd: func() error {
			return task1.Exec("uname")
		},
	})
	task2 := NewTask("Run uname command2", WithWriter(buf))
	task2.SetFuncs(ExecFuncParam{
		TargetCmd: nil,
		DepCmd:    nil,
		InstCmd: func() error {
			return task2.Exec("uname")
		},
	})
	job := NewJob([]Task{task1, task2})
	job.Run()

	expect := `[1/2 Run uname command1]
  => [exec] uname
  => Linux
=> Success install
[2/2 Run uname command2]
  => [exec] uname
  => Linux
=> Success install
`
	assert.Equal(t, expect, buf.String())
}
