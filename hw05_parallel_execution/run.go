package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorCounter struct {
	count int
	m     sync.RWMutex
}

func (e *ErrorCounter) Increment() {
	e.m.Lock()
	defer e.m.Unlock()
	e.count++
}

func (e *ErrorCounter) IsMax(max int) bool {
	if max < 0 {
		return false
	}
	e.m.RLock()
	defer e.m.RUnlock()
	return e.count >= max
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := new(sync.WaitGroup)
	tasksChannel := make(chan Task, n)
	var errCounter ErrorCounter

	// Start n goroutines to execute Tasks
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksChannel {
				if errCounter.IsMax(m) {
					continue
				}
				err := task()
				if err != nil {
					errCounter.Increment()
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
			tasksChannel <- task
		}
	}()

	wg.Wait()

	if errCounter.IsMax(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
