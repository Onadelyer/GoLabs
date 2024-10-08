package main

import (
	"fmt"
	"sync"
	"workerpool/workerpool"
)

func main() {
	var wg sync.WaitGroup
	counter := 0
	mutex := sync.Mutex{}

	workerPool := workerpool.NewWorkerPool(3, &wg)
	workerPool.Run()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		jobId := i + 1
		job := workerpool.Job{
			Id:          jobId,
			Description: fmt.Sprintf("Increment counter #%d", jobId),
			Run: func() {
				mutex.Lock()
				counter++
				fmt.Printf("Counter value: %d\n", counter)
				mutex.Unlock()
			},
		}
		workerPool.Jobs <- job
	}

	wg.Wait()
	close(workerPool.Jobs)
	fmt.Println("All jobs completed.")
}