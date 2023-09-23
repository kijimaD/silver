package silver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskInst(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	task := NewTask("Run uname command", WithWriter(buf))
	task.instCmd = func() error {
		err := task.Exec("uname")

		return err
	}
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
	task := NewTask("Run uname command", WithWriter(buf))
	task.depCmd = func() bool {
		return IsExistCmd("not_found_cmd")
	}
	task.instCmd = func() error {
		err := task.Exec("uname")

		return err
	}
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
	task := NewTask("Run uname command", WithWriter(buf))
	task.targetCmd = func() bool {
		return IsExistCmd("uname")
	}
	task.instCmd = func() error {
		err := task.Exec("uname")

		return err
	}
	task.Run()
	expect := `[Run uname command]
=> Already achieved, skip
`
	assert.Equal(t, expect, buf.String())
}
