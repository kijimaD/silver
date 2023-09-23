package silver

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type Task struct {
	name      string
	status    statusText
	targetCmd BoolFunc  // 条件。trueだと実行の必要がないとして、実行しない
	depCmd    BoolFunc  // 条件。trueだと依存関係を満たしているとして、実行する
	instCmd   ErrorFunc // 実行するコマンド
	w         io.Writer
}

type statusText string

const (
	waitExecuteST     = statusText("wait Execute")
	successInstallST  = statusText("Success install")
	failInstallST     = statusText("Fail install")
	notMetST          = statusText("Dependencies not met, skip")
	alreadyAchievedST = statusText("Already achieved, skip")
)

type BoolFunc func() bool
type ErrorFunc func() error

func NewTask(name string, w io.Writer) Task {
	t := Task{
		name:      name,
		status:    waitExecuteST,
		targetCmd: func() bool { return false },
		depCmd:    func() bool { return true },
		instCmd:   func() error { return nil },
		w:         w,
	}
	return t
}

func (t *Task) SetFuncs(target BoolFunc, dep BoolFunc, inst ErrorFunc) {
	t.targetCmd = target
	t.depCmd = dep
	t.instCmd = inst
}

func (t *Task) Run() {
	fmt.Fprintf(t.w, "[%s]\n", t.name)

	procs := []func() bool{
		t.processTarget,
		t.processDep,
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
	go t.displayOutput(stdout)
	go t.displayOutput(stderr)

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("コマンドの実行中にエラーが発生した%s", err)
	}

	return nil
}

func (t *Task) displayOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(t.w, "  => %s\n", scanner.Text())
	}
}

func (t *Task) processTarget() bool {
	ok := t.targetCmd()
	if ok {
		t.status = alreadyAchievedST
		return false
	}
	return true
}

func (t *Task) processDep() bool {
	ok := t.depCmd()
	if !ok {
		t.status = notMetST
		return false
	}
	return true
}
func (t *Task) processInst() bool {
	err := t.instCmd()
	if err != nil {
		t.status = failInstallST
		return false
	} else {
		t.status = successInstallST
	}
	return true
}
