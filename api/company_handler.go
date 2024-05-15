package api

type HotelHandler struct{}

func NewCompanyHandler() *HotelHandler {
	return &HotelHandler{}
}

func (h *HotelHandler) HandleGetCompanies() error {
	return nil
}

func (h *HotelHandler) HandleCreateCompany() error {
	return nil
}
