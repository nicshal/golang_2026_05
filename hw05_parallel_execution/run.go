package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	taskChan := make(chan Task)
	mu := sync.Mutex{}
	var errCount int

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				task, ok := <-taskChan
				if !ok {
					return
				}
				err := task()
				if err != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}
		}()
	}

	go func() {
		defer close(taskChan)
		for _, task := range tasks {
			mu.Lock()
			if errCount >= m && m > 0 {
				mu.Unlock()
				break
			}
			mu.Unlock()
			taskChan <- task
		}
	}()

	wg.Wait()

	if errCount >= m && m > 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
