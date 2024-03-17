package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// StatusResponse описывает структуру ответа на запрос статуса
type StatusResponse struct {
	ServiceName string    `json:"service_name"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to the Home Page!")
}

func StatusHandler(w http.ResponseWriter, _ *http.Request) {
	response := StatusResponse{
		ServiceName: "MyService",
		Status:      "OK",
		Timestamp:   time.Now(),
		Version:     "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
