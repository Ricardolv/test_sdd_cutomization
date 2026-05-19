package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/services"
)

type CustomerHandler struct {
	service *services.CustomerService
}

func NewCustomerHandler(service *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	c, err := h.service.Save(r.Context(), "", req.Name, req.Email, req.Phone, req.Notes)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toResponse(c))
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	c, err := h.service.Save(r.Context(), id, req.Name, req.Email, req.Phone, req.Notes)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toResponse(c))
}

func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	c, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toResponse(c))
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	customers, err := h.service.List(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}

	responses := make([]CustomerResponse, len(customers))
	for i, c := range customers {
		responses[i] = toResponse(&c)
	}

	writeJSON(w, http.StatusOK, responses)
}

func (h *CustomerHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Deactivate(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toResponse(c *models.Customer) CustomerResponse {
	return CustomerResponse{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Notes:     c.Notes,
		Active:    c.Active,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case models.ErrNameRequired:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrEmailRequired:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrEmailDuplicate:
		writeError(w, http.StatusConflict, "duplicate_email", err.Error())
	case models.ErrCustomerNotFound:
		writeError(w, http.StatusNotFound, "not_found", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected error occurred")
	}
}
