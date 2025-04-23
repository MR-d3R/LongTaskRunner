package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"longtaskrunner/api"
	"longtaskrunner/queue"
	"longtaskrunner/storage"
	"longtaskrunner/worker"
)

func main() {
	resultStorage := storage.NewInMemoryStorage()

	taskQueue := queue.NewInMemoryQueue()

	workerPool := worker.NewWorkerPool(10, taskQueue, resultStorage)

	apiServer := api.NewServer(taskQueue, resultStorage)

	workerPool.Start()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: apiServer.Router(),
	}

	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Обработка сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	workerPool.Stop()
	log.Println("Server stopped")
}
