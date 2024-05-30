package api

import (
	"bytes"
	"context"
	"encoding/json"
	"h-project/db"
	"h-project/internal/entity"
	"h-project/internal/file"
	"io/ioutil"
	"log/slog"
	"net/http"
	"time"
)

type CompanyHandler struct {
	store       *db.DB
	fileService file.Service
	logger      *slog.Logger
}

func NewCompanyHandler(store *db.DB, fileService file.Service, logger *slog.Logger) *CompanyHandler {
	return &CompanyHandler{
		store:       store,
		fileService: fileService,
		logger:      logger,
	}
}

func (h *CompanyHandler) HandleGetCompanies(w http.ResponseWriter, _ *http.Request) {
	companies, err := h.store.GetCompanies()
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companies)
}

func (h *CompanyHandler) HandleCreateCompany(w http.ResponseWriter, r *http.Request) {
	var company entity.Company

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal the request body into the Company struct
	err = json.Unmarshal(body, &company)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := company.Name + time.Now().String()
	reader := bytes.NewReader([]byte(name))
	dto := file.CreateFileDTO{
		Name:   name,
		Size:   int64(len(name)),
		Reader: reader,
	}

	//Передача данных в бакет
	ctx := context.Background()

	err = h.fileService.Create(ctx, file.BUCKET_NAME, dto)

	if err != nil {
		h.logger.Error(err.Error())
	}

	err = h.store.AddCompany(&company)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
