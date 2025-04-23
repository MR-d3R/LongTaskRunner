package model

import (
	"time"
)

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

type TaskRequest struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

type Task struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Params    map[string]interface{} `json:"params"`
	Status    string                 `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	StartedAt *time.Time             `json:"started_at,omitempty"`
	EndedAt   *time.Time             `json:"ended_at,omitempty"`
}

type TaskResult struct {
	TaskID    string                 `json:"task_id"`
	Status    string                 `json:"status"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	StartedAt *time.Time             `json:"started_at,omitempty"`
	EndedAt   *time.Time             `json:"ended_at,omitempty"`
}

type TaskResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
