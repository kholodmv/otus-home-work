package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	taskChan := make(chan Task, len(tasks))
	errorCount := 0

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				err := task()
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					if errorCount >= m {
						return
					}
				}
			}
		}()
	}

	wg.Wait()

	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
