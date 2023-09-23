package silver

import (
	"fmt"
)

type Task struct {
	name     string
	status   statusText
	instCmds []errorFunc
}

type statusText string

const (
	waitExecute     = statusText("wait Execute")
	successInstall  = statusText("Success install")
	failInstall     = statusText("Fail install")
	notMet          = statusText("Dependencies not met")
	alreadyAchieved = statusText("Already achieved")
)

type errorFunc func() error

func NewTask(name string, instCmds []errorFunc) Task {
	t := Task{
		name:     name,
		status:   waitExecute,
		instCmds: instCmds,
	}
	return t
}

func (t *Task) Run() {
	fmt.Printf("[%s]\n", t.name)

	for _, cmd := range t.instCmds {
		err := cmd()
		if err != nil {
			t.status = failInstall
		} else {
			t.status = successInstall
		}
	}

	fmt.Printf("=> %s\n", t.status)
}
