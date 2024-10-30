package workerpool

import (
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(3)

	if pool.WorkerCount() != 3 {
		t.Fatalf("Expected 3 workers, got %d", pool.WorkerCount())
	}

	for i := 1; i <= 10; i++ {
		pool.AddJob("job")
	}

	time.Sleep(5 * time.Second)

	if pool.WorkerCount() != 3 {
		t.Fatalf("Expected 3 workers, got %d", pool.WorkerCount())
	}

	pool.AddWorker()
	if pool.WorkerCount() != 4 {
		t.Fatalf("Expected 4 workers after adding, got %d", pool.WorkerCount())
	}

	if err := pool.RemoveWorker(2); err != nil {
		t.Fatal(err)
	}
	if pool.WorkerCount() != 3 {
		t.Fatalf("Expected 3 workers after removing, got %d", pool.WorkerCount())
	}

	if err := pool.RemoveWorker(10); err == nil {
		t.Fatal("Expected error for invalid worker ID")
	}

	if err := pool.RemoveRandomWorker(); err != nil {
		t.Fatal(err)
	}
	if pool.WorkerCount() != 2 {
		t.Fatalf("Expected 2 workers after random removal, got %d", pool.WorkerCount())
	}

	for i := 0; i < 2; i++ {
		pool.RemoveRandomWorker()
	}
	if err := pool.RemoveRandomWorker(); err == nil {
		t.Fatal("Expected error when trying to remove from empty pool")
	}

	pool.AddWorker() // Добавляем нового воркера
	if pool.WorkerCount() != 1 {
		t.Fatalf("Expected 1 worker after adding, got %d", pool.WorkerCount())
	}

	pool.Stop()
}

func TestRemoveWorkerErrors(t *testing.T) {
	pool := NewWorkerPool(3)

	if err := pool.RemoveWorker(0); err == nil {
		t.Fatal("Expected error for invalid worker ID (0)")
	}
	if err := pool.RemoveWorker(4); err == nil {
		t.Fatal("Expected error for invalid worker ID (4), which doesn't exist")
	}

	pool.Stop()
}

func TestRemoveRandomWorkerErrors(t *testing.T) {
	pool := NewWorkerPool(1)

	if err := pool.RemoveRandomWorker(); err != nil {
		t.Fatal(err)
	}

	if err := pool.RemoveRandomWorker(); err == nil {
		t.Fatal("Expected error when trying to remove from empty pool")
	}

	pool.Stop()
}

func TestConcurrentAddRemove(t *testing.T) {
	pool := NewWorkerPool(3)

	for i := 0; i < 5; i++ {
		go pool.AddWorker()
		go pool.RemoveRandomWorker()
	}

	time.Sleep(5 * time.Second)

	pool.Stop()
}

func TestMultipleJobs(t *testing.T) {
	pool := NewWorkerPool(3)

	for i := 1; i <= 50; i++ {
		pool.AddJob("job")
	}

	time.Sleep(10 * time.Second)

	pool.Stop()
}
