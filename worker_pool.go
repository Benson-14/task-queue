package main

import "sync"

type WorkerPool struct {
	workers []*Worker
	queue   *Queue
	wg      *sync.WaitGroup
	stop    chan struct{}
}

func NewWorkerPool(queue *Queue, size int) *WorkerPool {
	return &WorkerPool{
		queue:   queue,
		stop:    make(chan struct{}),
		wg:      &sync.WaitGroup{},
		workers: make([]*Worker, size),
	}
}

func (wp *WorkerPool) Start() {
	for i := range wp.workers {
		worker := NewWorker(i, wp.queue, wp.wg)
		wp.workers[i] = worker
		worker.Start()
	}
}

func (wp *WorkerPool) Stop() {
	for _, w := range wp.workers {
		w.Stop()
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
