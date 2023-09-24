package silver

// TODO: jobとtasksのbufを一元化したい

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
		task.SetStats(StatsWithIdx(i+1, taskCount))
		task.Run()
		j.Tasks[i] = task
	}

	for _, task := range j.Tasks {
		task.printStatus()
	}
}
