package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded  = errors.New("errors limit exceeded")
	ErrInvalidWorkersNumber = errors.New("invalid workers number")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrInvalidWorkersNumber
	}

	errMax := int32(m)
	errCount := new(int32)
	isExceededErrMax := func() bool {
		if errMax <= 0 {
			return false
		}
		return atomic.LoadInt32(errCount) >= errMax
	}

	wg := new(sync.WaitGroup)
	tasksChannel := make(chan Task)

	// Start n goroutines to execute Tasks
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksChannel {
				err := task()
				if err != nil {
					atomic.AddInt32(errCount, 1)
				}
			}
		}()
	}

	// Fill tasks channel
	wg.Add(1)
	go func() {
		defer close(tasksChannel)
		defer wg.Done()
		for _, task := range tasks {
			if isExceededErrMax() {
				break
			}
			tasksChannel <- task
		}
	}()

	wg.Wait()

	if isExceededErrMax() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
