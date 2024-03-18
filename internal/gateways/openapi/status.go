package openapi

import (
	"encoding/json"
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

/*
 * if /healthz path returns a success code, the kubelet considers the container to be alive and healthy
 * if the handler returns a failure code, the kubelet kills the container and restarts it
 */
func Status(serviceName, serviceVersion string) http.HandlerFunc {

	response := StatusResponse{
		ServiceName: serviceName,
		Status:      "OK",
		Timestamp:   time.Now(),
		Version:     serviceVersion,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
