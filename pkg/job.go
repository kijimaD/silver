package silver

type Job struct {
	Tasks []Task
}

func NewJob(tasks []Task) Job {
	j := Job{
		Tasks: tasks,
	}

	return j
}

func (j *Job) Run() {
	for _, t := range j.Tasks {
		t.Run()
	}
}
