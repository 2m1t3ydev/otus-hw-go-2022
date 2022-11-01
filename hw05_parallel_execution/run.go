package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorCounter struct {
	mutex    *sync.Mutex
	counter  int
	maxCount int
}

func (ec *ErrorCounter) Increment() {
	(*ec.mutex).Lock()
	ec.counter++
	(*ec.mutex).Unlock()
}

func (ec *ErrorCounter) GetValue() int {
	var value int

	(*ec.mutex).Lock()
	value = ec.counter
	(*ec.mutex).Unlock()

	return value
}

func (ec *ErrorCounter) ReachLimit() bool {
	if ec.maxCount <= 0 {
		return false
	}

	value := ec.GetValue()
	return value > ec.maxCount
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Consumer(c *chan Task, counter *ErrorCounter, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range *c {
		if !(*counter).ReachLimit() {
			res := task()
			if res != nil {
				(*counter).Increment()
			}
		}
	}
}

func Producer(tasks []Task, c *chan Task, counter *ErrorCounter) {
	defer close(*c)

	for _, task := range tasks {
		if (*counter).ReachLimit() {
			return
		}
		(*c) <- task
	}
}

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	c := make(chan Task)
	counter := ErrorCounter{
		mutex:    &mu,
		counter:  0,
		maxCount: m,
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go Consumer(&c, &counter, &wg)
	}

	Producer(tasks, &c, &counter)

	wg.Wait()

	if counter.ReachLimit() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
