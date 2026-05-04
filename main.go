package main

import (
	"encoding/json"
	"fmt"
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

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {

	// Email Task
	emailData := EmailPayload{
		To:      "ben@gmail.com",
		Subject: "greeting",
		Body:    "This is the body of the email",
	}

	task, err := NewTask("email", emailData)
	if err != nil {
		fmt.Println("Error creating task:", err)
		return
	}
	fmt.Println("Task Created")

	fmt.Println(task.ID)
	fmt.Println(task.CreatedAt)
	fmt.Println(task.Status)

	var emailPayload EmailPayload
	if err := task.UnmarshalPayload(&emailPayload); err != nil {
		fmt.Println("error unmarshalling", err)
	}
	fmt.Println("Unmarshal Payload :", emailPayload)
}
