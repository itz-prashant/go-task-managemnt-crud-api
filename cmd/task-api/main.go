package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itz-prashant/task-api/internal/config"
	"github.com/itz-prashant/task-api/internal/middleware"
	"github.com/itz-prashant/task-api/internal/storage/sqlite"
	"github.com/itz-prashant/task-api/internal/task"
)

func main() {
	cfg := config.MustLoad()

	log.Printf("[INFO] Configuration loaded for environment: %s", cfg.Env)

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize database: %v", err)
	}

	defer storage.Db.Close()
	log.Printf("[INFO] SQLite database connected at: %s", cfg.StoragePath)

	taskRepo := task.NewRepository(storage.Db)
	taskService := task.NewService(taskRepo)
	taskHandler := task.NewHandler(taskService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/tasks", taskHandler.HandleCreate)
	mux.HandleFunc("GET /api/v1/tasks", taskHandler.HandleList)
	mux.HandleFunc("GET /api/v1/tasks/{id}", taskHandler.HandleGet)
	mux.HandleFunc("PUT /api/v1/tasks/{id}", taskHandler.HandleUpdate)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", taskHandler.HandleDelete)

	log.Println("[INFO] All routes registered successfully")

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.Recoverer,
	)

	serverHandler := stack(mux)

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      serverHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	log.Printf("[INFO] Server configured and ready to start on %s", cfg.HTTPServer.Addr)

	go func() {
		log.Printf("[INFO] Server listening on http://localhost%s", cfg.HTTPServer.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Server listening failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("[INFO] Shutting down server gracefully...")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutDownCtx); err != nil {
		log.Fatalf("[INFO] Server forced to shutdown: %v", err)
	}
	log.Println("[INFO] Server exited cleanly. Goodbye!")
}
