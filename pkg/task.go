package silver

import (
	"fmt"
	"io"
	"os/exec"
)

type Task struct {
	name     string
	status   statusText
	depsCmds []boolFunc
	instCmds []errorFunc
	w        io.Writer
}

type statusText string

const (
	waitExecuteST     = statusText("wait Execute")
	successInstallST  = statusText("Success install")
	failInstallST     = statusText("Fail install")
	notMetST          = statusText("Dependencies not met, skip")
	alreadyAchievedST = statusText("Already achieved, skip")
)

type boolFunc func() bool
type errorFunc func() error

func NewTask(name string, w io.Writer) Task {
	t := Task{
		name:     name,
		status:   waitExecuteST,
		depsCmds: []boolFunc{},
		instCmds: []errorFunc{},
		w:        w,
	}
	return t
}

func (t *Task) Run() {
	fmt.Fprintf(t.w, "[%s]\n", t.name)

	procs := []func() bool{
		t.processDeps,
		t.processInst,
	}

	for _, proc := range procs {
		ok := proc()
		if !ok {
			break
		}
	}

	fmt.Fprintf(t.w, "=> %s\n", t.status)
}

func (t *Task) Exec(cmdtext string) error {
	fmt.Fprintf(t.w, "  => [exec] %s\n", cmdtext)

	cmd := exec.Command("bash", "-c", cmdtext)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("標準出力パイプ作成に失敗した%s", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("標準エラー出力パイプ作成に失敗した%s", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("コマンド開始に失敗した%s", err)
	}

	// リアルタイムに表示
	go displayOutput(stdout, t.w)
	go displayOutput(stderr, t.w)

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("コマンドの実行中にエラーが発生した%s", err)
	}

	return nil
}

func (t *Task) processDeps() bool {
	for _, cmd := range t.depsCmds {
		ok := cmd()
		if !ok {
			t.status = notMetST
			return false
		}
	}
	return true
}
func (t *Task) processInst() bool {
	for _, cmd := range t.instCmds {
		err := cmd()
		if err != nil {
			t.status = failInstallST
			return false
		} else {
			t.status = successInstallST
		}
	}
	return true
}
