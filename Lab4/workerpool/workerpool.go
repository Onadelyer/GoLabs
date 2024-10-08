package workerpool

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	Jobs       chan Job
	WorkersNum int
	Wg         *sync.WaitGroup
}

func NewWorkerPool(workersNum int, group *sync.WaitGroup) *WorkerPool {
	return &WorkerPool{
		Jobs:       make(chan Job, 10),
		WorkersNum: workersNum,
		Wg:         group,
	}
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.WorkersNum; i++ {
		go func(workerId int) {
			for job := range wp.Jobs {
				fmt.Printf("Worker %d starting job %d: %s\n", workerId, job.Id, job.Description)
				job.Run()
				fmt.Printf("Worker %d finished job %d\n", workerId, job.Id)
				wp.Wg.Done()
			}
		}(i + 1)
	}
}