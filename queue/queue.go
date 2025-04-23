package queue

import (
	"errors"
	"sync"

	"longtaskrunner/model"
)

// Интерфейс очереди задач
type TaskQueue interface {
	Push(task model.Task) error
	Pop() (model.Task, error)
	Size() int
}

// Очередь задач
type InMemoryQueue struct {
	tasks []model.Task
	mu    sync.Mutex
	cond  *sync.Cond
}

// Инициализация очереди задач
func NewInMemoryQueue() *InMemoryQueue {
	q := &InMemoryQueue{
		tasks: make([]model.Task, 0),
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Добавление новой задачи
func (q *InMemoryQueue) Push(task model.Task) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.tasks = append(q.tasks, task)
	q.cond.Signal() // Уведомляем ожидающий поп
	return nil
}

// Убираем задачу из очереди
func (q *InMemoryQueue) Pop() (model.Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Ждем если очередь пуста
	for len(q.tasks) == 0 {
		q.cond.Wait()
	}

	if len(q.tasks) == 0 {
		return model.Task{}, errors.New("queue is empty")
	}

	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task, nil
}

// Размер очереди
func (q *InMemoryQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.tasks)
}
