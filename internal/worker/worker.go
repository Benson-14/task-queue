package worker

import (
	"fmt"
	"sync"

	"github.com/Benson-14/task-queue/internal/queue"
	"github.com/Benson-14/task-queue/internal/task"
)

type Worker struct {
	id          int
	queue       *queue.Queue
	stop        chan struct{}
	wg          *sync.WaitGroup
	processed   []string
	mu          sync.Mutex
	statusStore *task.StatusStore
}

func NewWorker(id int, queue *queue.Queue, wg *sync.WaitGroup, store *task.StatusStore) *Worker {
	return &Worker{
		id:          id,
		queue:       queue,
		stop:        make(chan struct{}),
		wg:          wg,
		processed:   make([]string, 0),
		statusStore: store,
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		w.run()
	}()

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
		case t := <-w.queue.Tasks():
			t.Status = task.StatusRunning
			w.statusStore.Set(t)

			err := w.process(t)

			if err != nil {
				t.Status = task.StatusFailed
				t.Error = err.Error()
			} else {
				t.Status = task.StatusCompleted
			}
			w.statusStore.Set(t)
		}
	}
}

func (w *Worker) process(t *task.Task) error {
	w.mu.Lock()
	w.processed = append(w.processed, t.ID)
	w.mu.Unlock()
	fmt.Printf("Worker %d processed %s\n", w.id, t.ID)
	return nil
}

func (w *Worker) GetProcessed() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	result := make([]string, len(w.processed))
	copy(result, w.processed)
	return result
}
