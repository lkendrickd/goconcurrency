package semaphore

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task is a struct that holds the data for a task
// you could add other fields here like Data which you could use to pass
// data to the task to be processed.
type Task struct {
	ID          int
	IsProcessed bool
}

// WorkerWithSemaphore is a worker that processes tasks concurrently
// the semaphore which is a buffered channel is used to control the number of concurrent
// tasks that can be processed at a time.
type WorkerWithSemaphore struct {
	tasks []*Task
	// results is a channel that will hold the results of the work
	// this could be a string, struct or any other type that you want to return
	results       chan string
	semaphore     chan struct{}
	maxConcurrent int
}

// New creates a new WorkerWithSemaphore
func New(tasks []*Task, maxConcurrent int) *WorkerWithSemaphore {
	return &WorkerWithSemaphore{
		tasks: tasks,
		// buffer the channel with the number of tasks
		results: make(chan string, len(tasks)),
		// buffer the semaphore with the number of concurrent tasks allowed to run
		semaphore:     make(chan struct{}, maxConcurrent),
		maxConcurrent: maxConcurrent,
	}
}

// Work processes the tasks concurrently using anonymous goroutines
func (w *WorkerWithSemaphore) Work() {
	var wg sync.WaitGroup

	for _, task := range w.tasks {
		// acquire the semaphore this is achieved by writing to the channel
		// this will block if the channel is full. If full the code will wait until
		// a slot is available. I use an empty struct as the value because it
		// doesn't take up any memory.
		w.semaphore <- struct{}{}

		// once we can send to the semaphore we can increment the waitgroup
		wg.Add(1)

		// start a goroutine to process the task this will run concurrently for each task
		// that we are iterating over
		go func(task *Task) {
			defer wg.Done()
			defer func() { <-w.semaphore }()
			w.DoWork(task)
		}(task)
	}

	// after all the tasks are executed concurrently we wait for all the goroutines to finish
	// this means we wait for all the wg.Done() calls to be executed which will decrement the counter
	// to zero and this blocking call will proceed.
	wg.Wait()

	// close the results channel once all the tasks are processed
	// even though the channel is closed we can still read from it
	close(w.results)
}

// DoWork simulates work on a task
func (w *WorkerWithSemaphore) DoWork(task *Task) {
	if task.IsProcessed {
		return
	}

	workTime := time.Duration(rand.Intn(5)) * time.Second
	fmt.Printf("Task %d is running\n", task.ID)
	for i := 0; i < int(workTime.Seconds()); i++ {
		time.Sleep(1 * time.Second)
	}
	task.IsProcessed = true
	fmt.Printf("Task %d is Completed\n", task.ID)
	w.results <- fmt.Sprintf("Task %d completed in %v", task.ID, workTime)
}

// Results returns the results of the work
func (w *WorkerWithSemaphore) Results() <-chan string {
	return w.results
}
