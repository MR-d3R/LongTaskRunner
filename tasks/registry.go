package tasks

import (
	"sync"
)

// TaskHandler определяет интерфейс для обработчика задач
type TaskHandler interface {
	Execute(params map[string]interface{}) (map[string]interface{}, error)
}

// Registry - реестр типов задач и их обработчиков
type Registry struct {
	handlers map[string]TaskHandler
	mu       sync.RWMutex
}

func NewRegistry() Registry {
	registry := Registry{
		handlers: make(map[string]TaskHandler),
	}

	// Регистрируем обработчики по умолчанию
	registry.RegisterHandler("long_computation", NewLongComputationHandler())

	return registry
}

func (r *Registry) RegisterHandler(taskType string, handler TaskHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[taskType] = handler
}

func (r *Registry) GetHandler(taskType string) (TaskHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, exists := r.handlers[taskType]
	return handler, exists
}
