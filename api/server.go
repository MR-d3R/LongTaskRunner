package api

import (
	"encoding/json"
	"net/http"
	"time"

	"longtaskrunner/model"
	"longtaskrunner/queue"
	"longtaskrunner/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// API сервер
type Server struct {
	taskQueue     queue.TaskQueue
	resultStorage storage.ResultStorage
}

// Инициализация API сервера
func NewServer(taskQueue queue.TaskQueue, resultStorage storage.ResultStorage) *Server {
	return &Server{
		taskQueue:     taskQueue,
		resultStorage: resultStorage,
	}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/tasks", s.createTask)
		r.Get("/tasks/{taskID}", s.getTaskStatus)
		r.Get("/tasks/{taskID}/result", s.getTaskResult)
	})

	return r
}

// Создание задачи и добавление её в очередь
func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	var req model.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := model.Task{
		ID:        uuid.New().String(),
		Type:      req.Type,
		Params:    req.Params,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.taskQueue.Push(task); err != nil {
		http.Error(w, "Failed to queue task", http.StatusInternalServerError)
		return
	}

	resp := model.TaskResponse{
		ID:     task.ID,
		Status: task.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

// Получение статуса выполнения задачи
func (s *Server) getTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	result, found := s.resultStorage.Get(taskID)
	if !found {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	resp := model.TaskResponse{
		ID:     taskID,
		Status: result.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Получение результата выполнения задачи из хранилище
func (s *Server) getTaskResult(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	result, found := s.resultStorage.Get(taskID)
	if !found {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if result.Status != model.StatusCompleted {
		http.Error(w, "Task is not completed yet", http.StatusPreconditionFailed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
