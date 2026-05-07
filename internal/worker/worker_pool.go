package worker

import (
	"sync"

	"github.com/Benson-14/task-queue/internal/queue"
	"github.com/Benson-14/task-queue/internal/task"
)

type WorkerPool struct {
	workers     []*Worker
	queue       *queue.Queue
	wg          *sync.WaitGroup
	stop        chan struct{}
	statusStore *task.StatusStore
}

func NewWorkerPool(queue *queue.Queue, size int, store *task.StatusStore) *WorkerPool {
	return &WorkerPool{
		queue:       queue,
		stop:        make(chan struct{}),
		wg:          &sync.WaitGroup{},
		workers:     make([]*Worker, size),
		statusStore: store,
	}
}

func (wp *WorkerPool) Start() {
	for i := range wp.workers {
		worker := NewWorker(i, wp.queue, wp.wg, wp.statusStore)
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
