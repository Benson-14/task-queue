package api

import (
	"encoding/json"
	"net/http"

	"github.com/Benson-14/task-queue/internal/task"
	"github.com/Benson-14/task-queue/internal/utils"
)

type submitRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (srv *Server) handleSubmit(w http.ResponseWriter, r *http.Request) {
	var req submitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Type == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "task type is required")
		return
	}

	t, err := task.NewTask(req.Type, req.Payload)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	srv.statusStore.Set(t)

	if err := srv.queue.Push(t); err != nil {
		utils.RespondWithError(w, http.StatusServiceUnavailable, "queue is full, try again later")
		return
	}

	utils.RespondWithJSON(w, http.StatusAccepted, map[string]string{
		"id":     t.ID,
		"status": string(t.Status),
	})
}

func (srv *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "task id is required")
		return
	}

	t, ok := srv.statusStore.Get(id)
	if !ok {
		utils.RespondWithError(w, http.StatusNotFound, "task not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, t)
}

func (srv *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := srv.statusStore.Stats()
	utils.RespondWithJSON(w, http.StatusOK, stats)
}
