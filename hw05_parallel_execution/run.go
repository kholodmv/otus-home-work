package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var wg sync.WaitGroup

	taskChan := make(chan Task, len(tasks))
	errorChan := make(chan struct{}, m)

	var errorCount int32

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				curErrorCount := atomic.LoadInt32(&errorCount)
				if curErrorCount >= int32(m) {
					return
				}
				err := task()
				if err != nil {
					atomic.AddInt32(&errorCount, 1)
					if atomic.LoadInt32(&errorCount) >= int32(m) {
						errorChan <- struct{}{}
						return
					}
				}
			}
		}()
	}

	go func() {
		defer close(taskChan)
		for _, task := range tasks {
			if atomic.LoadInt32(&errorCount) >= int32(m) {
				return
			}
			taskChan <- task
		}
	}()

	wg.Wait()

	if atomic.LoadInt32(&errorCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
