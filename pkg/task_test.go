package silver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskInst(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	task := NewTask(
		"Run uname command",
		buf,
	)
	testfunc := func() error {
		err := task.Exec("uname")

		return err
	}

	task.instCmd = testfunc

	task.Run()

	expect := `[Run uname command]
  => [exec] uname
  => Linux
=> Success install
`
	assert.Equal(t, expect, buf.String())
}

// 依存関係unameがないので、実行しない。
func TestTaskNotMet(t *testing.T) {
	t.Parallel()
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

	task.depCmd = depsfunc
	task.instCmd = testfunc

	task.Run()

	expect := `[Run uname command]
=> Dependencies not met, skip
`
	assert.Equal(t, expect, buf.String())
}

// ターゲットのunameがすでにあるので実行しない。
func TestTaskAlreadyAchived(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	task := NewTask(
		"Run uname command",
		buf,
	)
	targetfunc := func() bool {
		return IsExistCmd("uname")
	}
	testfunc := func() error {
		err := task.Exec("uname")

		return err
	}

	task.targetCmd = targetfunc
	task.instCmd = testfunc

	task.Run()

	expect := `[Run uname command]
=> Already achieved, skip
`
	assert.Equal(t, expect, buf.String())
}
