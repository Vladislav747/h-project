package application

import (
	"h-project/api"
	"h-project/internal/gateways/openapi"
	"net/http"
)

func NewHTTPHandler(serviceName, serviceVersion string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", api.HomeHandler)
	mux.HandleFunc("/healthz", openapi.Healthz())
	mux.HandleFunc("/readyz", openapi.Readyz())
	mux.HandleFunc("/status", openapi.Status(serviceName, serviceVersion))

	return mux
}
