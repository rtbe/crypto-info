package job

import "errors"

//Job contains info about particular job
type Job []string

//Queue - job queue
type Queue []Job

//New returns new job
func New(ss ...string) Job {
	return Job(ss)
}

//NewQueue returns new jobs queue
func NewQueue() *Queue {
	return &Queue{}
}

//Push job into job queue
func (jq *Queue) Push(j Job) {
	*jq = append(*jq, j)
}

//Pop job from job queue
func (jq *Queue) Pop() (Job, error) {
	if len(*jq) == 0 {
		return nil, errors.New("You can not pop values from empty queue")
	}
	j := []Job(*jq)[0]
	*jq = []Job(*jq)[1:]
	return j, nil
}
