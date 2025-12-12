package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	config := LoadConfig()

	// Initialize storage
	storage := NewMemoryStorage(config.MaxHistorySize)
	defer storage.Close()

	// Initialize Ollama client
	ollamaClient := NewOllamaClient(config.OllamaURL, config.OllamaModel, config.OllamaTimeout, config.MaxTokens)

	// Initialize application
	app := &App{
		config:  config,
		storage: storage,
		ollama:  ollamaClient,
	}

	// Setup router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(loggingMiddleware)
	router.Use(rateLimitMiddleware(config.RateLimit))
	router.Use(securityHeadersMiddleware)

	// Register routes
	router.HandleFunc("/", app.indexHandler).Methods("GET")
	router.HandleFunc("/ask", app.askHandler).Methods("POST")
	router.HandleFunc("/history", app.historyHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Create server
	srv := &http.Server{
		Addr:         config.ServerAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on %s", config.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
