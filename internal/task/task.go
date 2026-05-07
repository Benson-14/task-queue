package task

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID        string
	Type      string
	Payload   json.RawMessage
	Status    TaskStatus
	CreatedAt time.Time
	Error     string
}

type Stats struct {
	Pending   int
	Running   int
	Completed int
	Failed    int
}

type StatusStore struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func NewStatusStore() *StatusStore {
	return &StatusStore{
		tasks: make(map[string]*Task),
	}
}

func (s *StatusStore) Set(task *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

func (s *StatusStore) Get(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

func (s *StatusStore) ListByStatus(status TaskStatus) []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Task
	for _, task := range s.tasks {
		if task.Status == status {
			result = append(result, task)
		}
	}
	return result
}

func (s *StatusStore) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := Stats{}
	for _, task := range s.tasks {
		switch task.Status {
		case StatusPending:
			stats.Pending++
		case StatusRunning:
			stats.Running++
		case StatusCompleted:
			stats.Completed++
		case StatusFailed:
			stats.Failed++
		}
	}
	return stats
}

func NewTask(taskType string, payload any) (*Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &Task{
		ID:        generateID(),
		Type:      taskType,
		Payload:   data,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		Error:     "",
	}, nil
}

func generateID() string {
	return fmt.Sprintf("task-%d", time.Now().UnixNano())
}

func (t *Task) UnmarshalPayload(v any) error {
	return json.Unmarshal(t.Payload, v)
}
