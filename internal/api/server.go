package api

import (
	"net/http"

	"github.com/Benson-14/task-queue/internal/queue"
	"github.com/Benson-14/task-queue/internal/task"
)

type Server struct {
	router      *http.ServeMux
	queue       *queue.Queue
	statusStore *task.StatusStore
}

func NewServer(q *queue.Queue, store *task.StatusStore) *Server {
	srv := &Server{
		router:      http.NewServeMux(),
		queue:       q,
		statusStore: store,
	}
	srv.routes()
	return srv
}

func (srv *Server) routes() {
	srv.router.HandleFunc("GET /health", srv.healthHandler)
	srv.router.HandleFunc("POST /tasks", srv.handleSubmit)
	srv.router.HandleFunc("GET /tasks/stats", srv.handleStats)
	srv.router.HandleFunc("GET /tasks/{id}", srv.handleStatus)
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
