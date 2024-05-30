package application

import (
	"h-project/api"
	"h-project/db"
	"h-project/internal/file"
	"h-project/internal/gateways/openapi"
	"log/slog"
	"net/http"
)

func NewHTTPHandler(serviceName, serviceVersion string, db *db.DB, fileService file.Service, logger *slog.Logger) http.Handler {

	companyHandler := api.NewCompanyHandler(db, fileService, logger)
	fileHandler := file.NewFileHandler(fileService, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", openapi.Healthz())
	mux.HandleFunc("/readyz", openapi.Readyz())
	mux.HandleFunc("/status", openapi.Status(serviceName, serviceVersion))
	mux.HandleFunc("/companies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			companyHandler.HandleGetCompanies(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/companies/import", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			companyHandler.HandleCreateCompany(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/file/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fileHandler.GetFile(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			fileHandler.DeleteFile(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/file/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fileHandler.CreateFile(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
