package api

import (
	"encoding/json"
	"h-project/db"
	"h-project/internal/entity"
	"io/ioutil"
	"log/slog"
	"net/http"
)

type CompanyHandler struct {
	store  *db.DB
	logger *slog.Logger
}

func NewCompanyHandler(store *db.DB, logger *slog.Logger) *CompanyHandler {
	return &CompanyHandler{
		store:  store,
		logger: logger,
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

	err = h.store.AddCompany(&company)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
