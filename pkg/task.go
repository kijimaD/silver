package silver

import "fmt"

type statusText string

const (
	waitExecute     = statusText("wait Execute")
	successInstall  = statusText("Success install")
	failInstall     = statusText("Fail install")
	notMet          = statusText("Dependencies not met")
	alreadyAchieved = statusText("Already achieved")
)

type Task struct {
	name     string
	status   statusText
	instCmds []func()
}

func NewTask(name string, instCmds []func()) Task {
	t := Task{
		name:     name,
		status:   waitExecute,
		instCmds: instCmds,
	}
	return t
}

func (t *Task) Run() {
	fmt.Println("実行開始:", t.name)

	for _, cmd := range t.instCmds {
		cmd()
	}
}
