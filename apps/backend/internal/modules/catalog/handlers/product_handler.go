package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/repositories"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/services"
)

type ProductHandler struct {
	saveUseCase   *services.SaveProductUseCase
	deleteUseCase *services.DeleteProductUseCase
	repo          repositories.ProductRepository
}

func NewProductHandler(saveUseCase *services.SaveProductUseCase, deleteUseCase *services.DeleteProductUseCase, repo repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{
		saveUseCase:   saveUseCase,
		deleteUseCase: deleteUseCase,
		repo:          repo,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	input := services.SaveProductInput{
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Status:          models.ProductStatus(req.Status),
		AvailableOnline: req.AvailableOnline,
		Featured:        req.Featured,
		AllowsPreOrder:  req.AllowsPreOrder,
	}

	if err := h.saveUseCase.Execute(r.Context(), input); err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "product created"})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	input := services.SaveProductInput{
		ID:              id,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Status:          models.ProductStatus(req.Status),
		AvailableOnline: req.AvailableOnline,
		Featured:        req.Featured,
		AllowsPreOrder:  req.AllowsPreOrder,
	}

	if err := h.saveUseCase.Execute(r.Context(), input); err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "product updated"})
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toResponse(p))
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, total, err := h.repo.List(r.Context(), offset, limit)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = toResponse(&p)
	}

	writeJSON(w, http.StatusOK, ListResponse{
		Items: responses,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toResponse(p *models.Product) ProductResponse {
	return ProductResponse{
		ID:              p.ID,
		Name:            p.Name,
		Description:     p.Description,
		Price:           p.Price,
		Status:          string(p.Status),
		AvailableOnline: p.AvailableOnline,
		Featured:        p.Featured,
		AllowsPreOrder:  p.AllowsPreOrder,
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
	case models.ErrNameTooShort:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrNameTooLong:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrDescriptionTooLong:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrPriceRequired:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrPriceNegative:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrPricePrecision:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrStatusRequired:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrStatusInvalid:
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case models.ErrProductNotFound:
		writeError(w, http.StatusNotFound, "not_found", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected error occurred")
	}
}
