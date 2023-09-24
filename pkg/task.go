package silver

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Task struct {
	name     string
	status   statusText
	execFunc execFunc
	w        io.Writer
	Stats    Stats
}

type execFunc struct {
	targetCmd BoolFunc  // 条件。trueだと実行の必要がないとして、実行しない
	depCmd    BoolFunc  // 条件。trueだと依存関係を満たしているとして、実行する
	instCmd   ErrorFunc // 実行するコマンド
}

type ExecFuncParam struct {
	TargetCmd BoolFunc
	DepCmd    BoolFunc
	InstCmd   ErrorFunc
}

type Stats struct {
	CurrentIdx int
	AllLen     int
}

type statusText string

const (
	waitExecuteST     = statusText("wait Execute")
	successInstallST  = statusText("Success install")
	failInstallST     = statusText("Fail install")
	notMetST          = statusText("Dependencies not met, skip")
	alreadyAchievedST = statusText("Already achieved, skip")
)

type (
	BoolFunc  func() bool
	ErrorFunc func() error
)

func NewTask(name string, options ...TaskOption) Task {
	ef := execFunc{
		targetCmd: func() bool { return false },
		depCmd:    func() bool { return true },
		instCmd:   func() error { return nil },
	}
	s := Stats{
		CurrentIdx: 0,
		AllLen:     0,
	}
	t := Task{
		name:     name,
		status:   waitExecuteST,
		execFunc: ef,
		Stats:    s,
		w:        os.Stdout,
	}
	for _, option := range options {
		option(&t)
	}

	return t
}

func (t *Task) SetFuncs(execFuncParam ExecFuncParam) {
	if execFuncParam.TargetCmd != nil {
		t.execFunc.targetCmd = execFuncParam.TargetCmd
	}
	if execFuncParam.DepCmd != nil {
		t.execFunc.depCmd = execFuncParam.DepCmd
	}
	if execFuncParam.InstCmd != nil {
		t.execFunc.instCmd = execFuncParam.InstCmd
	}
}

func (t *Task) Run() {
	fmt.Fprintf(t.w, "[%d/%d %s]\n", t.Stats.CurrentIdx, t.Stats.AllLen, t.name)

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
		return fmt.Errorf("標準出力パイプ作成に失敗した%w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("標準エラー出力パイプ作成に失敗した%w", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("コマンド開始に失敗した%w", err)
	}

	// リアルタイムに表示
	go t.displayOutput(stdout)
	go t.displayOutput(stderr)

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("コマンドの実行中にエラーが発生した%w", err)
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
	ok := t.execFunc.targetCmd()
	if ok {
		t.status = alreadyAchievedST

		return false
	}

	return true
}

func (t *Task) processDep() bool {
	ok := t.execFunc.depCmd()
	if !ok {
		t.status = notMetST

		return false
	}

	return true
}

func (t *Task) processInst() bool {
	err := t.execFunc.instCmd()
	if err != nil {
		t.status = failInstallST

		return false
	}
	t.status = successInstallST

	return true
}

func (t *Task) SetStats(options ...StatsOption) {
	for _, option := range options {
		option(&t.Stats)
	}
}
