package workerpool

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Worker struct {
	id      int
	jobChan chan string
	done    chan struct{}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-w.jobChan:
			if !ok {
				return
			}
			fmt.Printf("Worker %d processing job: %s\n", w.id, job)
			time.Sleep(time.Second)
		case <-w.done:
			return
		}
	}
}

type WorkerPool struct {
	workers  []*Worker
	jobChan  chan string
	mu       sync.Mutex
	workerWg sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		workers: make([]*Worker, 0, numWorkers),
		jobChan: make(chan string),
	}
	for i := 0; i < numWorkers; i++ {
		pool.AddWorker()
	}
	return pool
}

func (p *WorkerPool) AddWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	worker := &Worker{
		id:      len(p.workers) + 1,
		jobChan: p.jobChan,
		done:    make(chan struct{}),
	}
	p.workers = append(p.workers, worker)
	p.workerWg.Add(1)
	go worker.Start(&p.workerWg)
}

func (p *WorkerPool) RemoveWorker(id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if id < 1 || id > len(p.workers) {
		return fmt.Errorf("invalid worker ID: %d", id)
	}

	worker := p.workers[id-1]
	close(worker.done)
	p.workers = append(p.workers[:id-1], p.workers[id:]...)
	return nil
}

func (p *WorkerPool) RemoveRandomWorker() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.workers) == 0 {
		return fmt.Errorf("no workers to remove")
	}

	randomIndex := rand.Intn(len(p.workers))
	worker := p.workers[randomIndex]
	close(worker.done)
	p.workers = append(p.workers[:randomIndex], p.workers[randomIndex+1:]...)
	return nil
}

func (p *WorkerPool) AddJob(job string) {
	p.jobChan <- job
}

func (p *WorkerPool) Stop() {
	close(p.jobChan)
	p.workerWg.Wait()
}

func (p *WorkerPool) WorkerCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.workers)
}
