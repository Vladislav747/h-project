package application

import (
	"h-project/api"
	"h-project/internal/gateways/openapi"
	"net/http"
)

func NewHTTPHandler(serviceName, serviceVersion string) http.Handler {

	companyHandler := api.NewCompanyHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("/", api.HomeHandler)
	mux.HandleFunc("/healthz", openapi.Healthz())
	mux.HandleFunc("/readyz", openapi.Readyz())
	mux.HandleFunc("/status", openapi.Status(serviceName, serviceVersion))
	mux.HandleFunc("/companies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			companyHandler.HandleGetCompanies()
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
