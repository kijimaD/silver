package silver

import (
	"fmt"
)

type Task struct {
	name     string
	status   statusText
	depsCmds []boolFunc
	instCmds []errorFunc
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

func NewTask(name string, depsCmds []boolFunc, instCmds []errorFunc) Task {
	t := Task{
		name:     name,
		status:   waitExecuteST,
		depsCmds: depsCmds,
		instCmds: instCmds,
	}
	return t
}

func (t *Task) Run() {
	fmt.Printf("[%s]\n", t.name)

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

	fmt.Printf("=> %s\n", t.status)
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
