package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type counter struct {
	counter int
	mu      sync.Mutex
}

func newCounter() *counter {
	return &counter{
		counter: 0,
		mu:      sync.Mutex{},
	}
}

func (c *counter) increase() {
	c.mu.Lock()
	c.counter++
	c.mu.Unlock()
}

func (c *counter) get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counter
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.
	wg := new(sync.WaitGroup)
	taskChan := make(chan Task)
	errorsCounter := newCounter()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskChan {
				err := t()
				if err != nil {
					errorsCounter.increase()
				}
			}
		}()
	}

	for _, task := range tasks {
		taskChan <- task

		if errorsCounter.get() >= m {
			close(taskChan)
			wg.Wait()
			return ErrErrorsLimitExceeded
		}
	}

	close(taskChan)
	wg.Wait()
	return nil
}
