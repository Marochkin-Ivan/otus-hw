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
	// Place your code here.
	wg := new(sync.WaitGroup)
	taskChan := make(chan Task)
	var errorsCounter atomic.Uint64

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskChan {
				err := t()
				if err != nil {
					errorsCounter.Add(1)
				}
			}
		}()
	}

	for _, task := range tasks {
		taskChan <- task

		if errorsCounter.Load() >= uint64(m) {
			close(taskChan)
			wg.Wait()
			return ErrErrorsLimitExceeded
		}
	}

	close(taskChan)
	wg.Wait()
	return nil
}
