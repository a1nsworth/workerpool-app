package main

import (
	"fmt"
	"time"

	"vk/workerpool"
)

func main() {
	workerCount := 10
	pool := workerpool.NewWorkerPool(workerCount)

	// Добавляем задачи в пул
	for i := 1; i <= 10; i++ {
		pool.AddJob(fmt.Sprintf("job %d", i))
	}

	time.Sleep(5 * time.Second)

	fmt.Println("add new worker")
	pool.AddWorker()

	for i := 1; i <= 11; i++ {
		pool.AddJob(fmt.Sprintf("job %d", i))
	}

	time.Sleep(10 * time.Second)

	pool.Stop()
}
