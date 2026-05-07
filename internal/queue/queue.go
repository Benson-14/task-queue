package queue

import (
	"errors"
	"time"

	"github.com/Benson-14/task-queue/internal/task"
)

// Queue is a channel-based task queue
type Queue struct {
	tasks chan *task.Task
	stop  chan struct{}
}

// NewQueue creates a queue with the given capacity
func NewQueue(capacity int) *Queue {
	return &Queue{
		tasks: make(chan *task.Task, capacity),
		stop:  make(chan struct{}),
	}
}

// Push adds a task to the queue (non-blocking)
// Returns error if queue is full
func (q *Queue) Push(task *task.Task) error {
	select {
	case q.tasks <- task:
		return nil
	default:
		return errors.New("queue is full")
	}
}

// Pop removes and returns a task (non-blocking)
// Returns (nil, false) if queue is empty
func (q *Queue) Pop() (*task.Task, bool) {
	select {
	case task := <-q.tasks:
		return task, true
	default:
		return nil, false
	}
}

// PopWait blocks until a task is available or timeout
func (q *Queue) PopWait(timeout time.Duration) (*task.Task, bool) {
	select {
	case task := <-q.tasks:
		return task, true
	case <-time.After(timeout):
		return nil, false
	}
}

// Len returns the current number of tasks in queue
func (q *Queue) Len() int {
	return len(q.tasks)
}

// Tasks returns the underlying task channel for workers to receive from
func (q *Queue) Tasks() <-chan *task.Task {
	return q.tasks
}

// Close signals the queue to stop
func (q *Queue) Close() {
	close(q.stop)
}
