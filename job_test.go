package silver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobRun(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	task1 := NewTask("Run uname command1", TaskWithWriter(buf))
	task1.SetFuncs(ExecFuncParam{
		TargetCmd: nil,
		DepCmd:    nil,
		InstCmd: func() error {
			return task1.Exec("uname")
		},
	})
	task2 := NewTask("Run uname command2", TaskWithWriter(buf))
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
  => [result] Success execute
[2/2 Run uname command2]
  => [exec] uname
  => Linux
  => [result] Success execute
[1/2 Run uname command1] Success execute
[2/2 Run uname command2] Success execute
`
	assert.Equal(t, expect, buf.String())
}

func TestJobRunMulti(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	task1 := NewTask("Run1", TaskWithWriter(buf))
	task1.SetFuncs(ExecFuncParam{
		TargetCmd: nil,
		DepCmd:    nil,
		InstCmd: func() error {
			return task1.Exec("echo test1 && echo test2 && echo test3")
		},
	})
	task2 := NewTask("Run2", TaskWithWriter(buf))
	task2.SetFuncs(ExecFuncParam{
		TargetCmd: nil,
		DepCmd:    nil,
		InstCmd: func() error {
			return task2.Exec("echo test1 && echo test2")
		},
	})
	job := NewJob([]Task{task1, task2})
	job.Run()

	expect := `[1/2 Run1]
  => [exec] echo test1 && echo test2 && echo test3
  => test1
  => test2
  => test3
  => [result] Success execute
[2/2 Run2]
  => [exec] echo test1 && echo test2
  => test1
  => test2
  => [result] Success execute
[1/2 Run1] Success execute
[2/2 Run2] Success execute
`
	assert.Equal(t, expect, buf.String())
}
