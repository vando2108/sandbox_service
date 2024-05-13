package workqueue

import (
	"fmt"
	"sync"
	"time"
)

type Task func() error

type WorkQueue struct {
	tasks      chan Task
	mutex      sync.Mutex
	ticker     *time.Ticker
	bucketSize int
}

func NewWorkQueue(bufferSize int, interval time.Duration, bucketSize int) *WorkQueue {
	return &WorkQueue{
		tasks:      make(chan Task, bufferSize),
		mutex:      sync.Mutex{},
		ticker:     time.NewTicker(interval),
		bucketSize: bucketSize,
	}
}

func (wq *WorkQueue) AddTask(task Task) {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	wq.tasks <- task
}

func (wq *WorkQueue) Start() {
	go func() {
		for range wq.ticker.C {
			wq.processTasks()
		}
	}()
}

func (wq *WorkQueue) Stop() {
	wq.ticker.Stop()
}

func (wq *WorkQueue) processTasks() error {
	for i := 0; i < wq.bucketSize; i++ {
		task, ok := <-wq.tasks
		if !ok {
			return fmt.Errorf("Channel has closed")
		}

		go task()
	}

	return nil
}
