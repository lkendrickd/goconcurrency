package semaphore

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a task to be processed
// using a pointer to a bool to reflect changes outside of DoWork
type Task struct {
	ID          int
	IsProcessed *bool
}

// WorkerWithSemaphore is a worker that processes tasks concurrently
// the semaphore is used to limit the number of concurrent tasks
// the maxConcurrent parameter specifies the maximum number of tasks
// that can ever run concurrently.
type WorkerWithSemaphore struct {
	tasks         []*Task
	results       chan string
	semaphore     chan struct{}
	maxConcurrent int
}

// NewWorkerWithSemaphore creates a new concreate instance of WorkerWithSemaphore
func NewWorkerWithSemaphore(tasks []*Task, maxConcurrent int) *WorkerWithSemaphore {
	return &WorkerWithSemaphore{
		tasks:   tasks,
		results: make(chan string, len(tasks)),
		// create a buffered channel to ratelimit on the maxConcurrent value
		semaphore:     make(chan struct{}, maxConcurrent),
		maxConcurrent: maxConcurrent,
	}
}

// Work processes the tasks concurrently using anonymous goroutiness
func (w *WorkerWithSemaphore) Work() {
	var wg sync.WaitGroup

	for _, task := range w.tasks {
		// Acquire the semaphore
		w.semaphore <- struct{}{}
		wg.Add(1)

		go func(task *Task) {
			defer wg.Done()
			defer func() { <-w.semaphore }()
			w.DoWork(task)
		}(task)
	}

	// Wait for all tasks to complete and close the results channel
	go func() {
		wg.Wait()
		close(w.results)
	}()
}

// DoWork simulates work by sleeping for a random duration
func (w *WorkerWithSemaphore) DoWork(task *Task) {
	if *task.IsProcessed {
		return
	}

	workTime := time.Duration(rand.Intn(5)) * time.Second
	fmt.Printf("Task %d is running\n", task.ID)
	for i := 0; i < int(workTime.Seconds()); i++ {
		time.Sleep(1 * time.Second)
	}

	*task.IsProcessed = true
	fmt.Printf("Task %d is Completed\n", task.ID)
	w.results <- fmt.Sprintf("Task %d completed in %v", task.ID, workTime)
}

// Results returns a channel to receive results from the worker
func (w *WorkerWithSemaphore) Results() <-chan string {
	return w.results
}
