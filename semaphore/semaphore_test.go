package semaphore_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lkendrickd/concurrency/semaphore"
)

func TestWorkerWithSemaphore(t *testing.T) {
	tests := []struct {
		name          string
		taskCount     int
		maxConcurrent int
		expected      int
	}{
		{
			name:          "single task",
			taskCount:     1,
			maxConcurrent: 1,
			expected:      1,
		},
		{
			name:          "multiple tasks, under limit",
			taskCount:     2,
			maxConcurrent: 4,
			expected:      2,
		},
		{
			name:          "multiple tasks, over limit",
			taskCount:     4,
			maxConcurrent: 2,
			expected:      4,
		},
		{
			name:          "multiple tasks, at limit",
			taskCount:     4,
			maxConcurrent: 4,
			expected:      4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks := make([]*semaphore.Task, tt.taskCount)
			for i := range tasks {
				var isProcessed bool
				tasks[i] = &semaphore.Task{
					ID:          i + 1,
					IsProcessed: isProcessed,
				}
			}

			worker := semaphore.New(tasks, tt.maxConcurrent)
			startTime := time.Now()
			worker.Work()

			results := make([]string, 0)
			for result := range worker.Results() {
				results = append(results, result)
			}

			if got := len(results); got != tt.expected {
				t.Errorf("expected %d results, got %d", tt.expected, got)
			}

			if len(results) != tt.taskCount {
				t.Errorf("expected %d tasks to complete, got %d", tt.taskCount, len(results))
			}

			for _, task := range tasks {
				if !task.IsProcessed {
					t.Errorf("Task %d was not processed", task.ID)
				}
			}

			fmt.Printf("Test '%s' completed in %v\n", tt.name, time.Since(startTime))
		})
	}
}
