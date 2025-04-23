package storage

import (
	"sync"

	"longtaskrunner/model"
)

// Интерфейся управления хранилищем
type ResultStorage interface {
	Set(taskID string, result model.TaskResult) error
	Get(taskID string) (model.TaskResult, bool)
}

// Хранилище результатов
type InMemoryStorage struct {
	results map[string]model.TaskResult
	mu      sync.RWMutex
}

// Инициализация хранилища результатов
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		results: make(map[string]model.TaskResult),
	}
}

// Сохранение результата в хранилище
func (s *InMemoryStorage) Set(taskID string, result model.TaskResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[taskID] = result
	return nil
}

// Получение результата из хранилища
func (s *InMemoryStorage) Get(taskID string) (model.TaskResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, found := s.results[taskID]
	return result, found
}
