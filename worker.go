package main

import (
	"fmt"
	"sync"
)

type Worker struct {
	id        int
	queue     *Queue
	stop      chan struct{}
	wg        *sync.WaitGroup
	processed []string
	mu        sync.Mutex
}

func NewWorker(id int, queue *Queue, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:        id,
		queue:     queue,
		stop:      make(chan struct{}),
		wg:        wg,
		processed: make([]string, 0),
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)

	w.wg.Go(func() {
		defer w.wg.Done()
		w.run()
	})
}

func (w *Worker) Stop() {
	close(w.stop)
}

func (w *Worker) Wait() {
	w.wg.Wait()
}

func (w *Worker) run() {
	for {
		select {
		case <-w.stop:
			return
		case task := <-w.queue.tasks:
			w.process(task)
		}
	}
}

func (w *Worker) process(task *Task) {
	w.mu.Lock()
	w.processed = append(w.processed, task.ID)
	w.mu.Unlock()
	fmt.Printf("Worker %d processed %s\n", w.id, task.ID)
}

func (w *Worker) GetProcessed() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	result := make([]string, len(w.processed))
	copy(result, w.processed)
	return result
}
