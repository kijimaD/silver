package silver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	tsize "github.com/kopoli/go-terminal-size"
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
	StartedAt  time.Time
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
		StartedAt:  time.Now(),
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

	fmt.Fprintf(t.w, "  [result]=> %s\n", t.status)
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

// 実行結果をタイマーとともに出力する。標準出力と標準エラー出力で使っている
// FIXME: 出力がすべて出ないことがある。sleepを最後に入れると出るので、出力する前にループが終了しているのだろう
// 例) $ echo hello && sleep 2 && echo hello && echo hello
// => hello だけ
func (t *Task) displayOutput(r io.Reader) {
	const timerDisplayPrecision = 1 // `1.1` 表示秒数の小数精度
	const secondDisplayLen = 1      // `s` 秒数の単位文字列の長さ
	scanner := bufio.NewScanner(r)
	done := make(chan bool)

	s, err := tsize.GetSize() // current terminal size
	if err != nil {
		log.Fatal(err)
	}

	// 経過時間を書き換えて表示する
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				scannedText := scanner.Text()
				scannedText = strings.ReplaceAll(scannedText, " ", "")

				if len(scannedText) > 0 {
					diff := time.Now().Sub(t.Stats.StartedAt)
					head := fmt.Sprintf("  => %s", scannedText)
					timer := fmt.Sprintf(
						"%*.*fs",
						s.Width-len(head)-secondDisplayLen, // 最後の秒数より右にある"s"の分の1文字を引く
						timerDisplayPrecision,
						diff.Seconds(),
					)
					fmt.Fprintf(t.w, "\r%s%s", head, timer)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// 行を次に進める
	for scanner.Scan() {
		fmt.Fprintf(t.w, "\n")
		time.Sleep(100 * time.Millisecond)
	}
	done <- true
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

func (t *Task) printStatus() {
	fmt.Fprintf(t.w, "[%d/%d %s] %s\n", t.Stats.CurrentIdx, t.Stats.AllLen, t.name, t.status)
}
