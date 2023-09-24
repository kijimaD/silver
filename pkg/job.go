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
	taskCount := len(j.Tasks)
	for i, task := range j.Tasks {
		task.Stats.SetStats(i+1, taskCount)
		task.Run()
	}
}
