package main

import (
	"context"
	"fmt"
	"h-project/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	port := os.Getenv("APPLICATION_PORT")
	if port == "" {
		fmt.Println("APPLICATION_PORT is not set. Using default port :8080")
		port = ":8080" // Значение порта по умолчанию
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", api.HomeHandler)
	mux.HandleFunc("/status", api.StatusHandler)

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Printf("Server is running on http://localhost%s\n", port)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = <-stopChan
		fmt.Println("Gracefully shutting down...")
		_ = server.Shutdown(ctx)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Error starting server: %s\n", err)
	}

	fmt.Println("App Starting")
}
