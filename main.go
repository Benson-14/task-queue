package main

import (
	"log"
	"net/http"

	"github.com/Benson-14/task-queue/internal/api"
	"github.com/Benson-14/task-queue/internal/queue"
	"github.com/Benson-14/task-queue/internal/task"
	"github.com/Benson-14/task-queue/internal/worker"
)

func main() {

	q := queue.NewQueue(100)
	store := task.NewStatusStore()

	pool := worker.NewWorkerPool(q, 3)
	pool.Start()
	defer pool.Stop()

	srv := api.NewServer(q, store)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
