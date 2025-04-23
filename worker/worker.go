package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"longtaskrunner/model"
	"longtaskrunner/queue"
	"longtaskrunner/storage"
	"longtaskrunner/tasks"
)

// ПУл воркеров
type WorkerPool struct {
	numWorkers    int
	taskQueue     queue.TaskQueue
	resultStorage storage.ResultStorage
	registry      tasks.Registry
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
}

// Создание пула воркеров
func NewWorkerPool(numWorkers int, taskQueue queue.TaskQueue, resultStorage storage.ResultStorage) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		numWorkers:    numWorkers,
		taskQueue:     taskQueue,
		resultStorage: resultStorage,
		registry:      tasks.NewRegistry(),
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *WorkerPool) Stop() {
	p.cancel()
	p.wg.Wait()
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case <-p.ctx.Done():
			log.Printf("Worker %d stopped", id)
			return
		default:
			task, err := p.taskQueue.Pop()
			if err != nil {
				log.Printf("Worker %d: Error popping task: %v", id, err)
				continue
			}

			p.processTask(task)
		}
	}
}

func (p *WorkerPool) processTask(task model.Task) {
	log.Printf("Processing task %s of type %s", task.ID, task.Type)

	now := time.Now()
	task.Status = model.StatusProcessing
	task.StartedAt = &now

	// Сохраняем начальный статус
	result := model.TaskResult{
		TaskID:    task.ID,
		Status:    task.Status,
		CreatedAt: task.CreatedAt,
		StartedAt: task.StartedAt,
	}
	p.resultStorage.Set(task.ID, result)

	// Получаем обработчик для типа задачи
	handler, exists := p.registry.GetHandler(task.Type)
	if !exists {
		endTime := time.Now()
		result.Status = model.StatusFailed
		result.Error = "Unknown task type"
		result.EndedAt = &endTime
		p.resultStorage.Set(task.ID, result)
		log.Printf("Task %s failed: unknown type %s", task.ID, task.Type)
		return
	}

	// Выполняем задачу
	taskResult, err := handler.Execute(task.Params)

	// Обновляем результат
	endTime := time.Now()
	result.EndedAt = &endTime

	if err != nil {
		result.Status = model.StatusFailed
		result.Error = err.Error()
	} else {
		result.Status = model.StatusCompleted
		result.Result = taskResult
	}

	// Сохраняем результат
	p.resultStorage.Set(task.ID, result)

	log.Printf("Task %s processed with status %s", task.ID, result.Status)
}
