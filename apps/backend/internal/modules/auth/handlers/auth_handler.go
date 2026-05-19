package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	u, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		handleAuthServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
	})
}

func handleAuthServiceError(w http.ResponseWriter, err error) {
	switch {
	case models.IsValidationError(err):
		errs := models.GetValidationErrors(err)
		codes := make([]string, len(errs))
		for i, e := range errs {
			codes[i] = mapErrorToCode(e)
		}
		writeValidationErrors(w, http.StatusUnprocessableEntity, "validation_error", codes)
	case err == models.ErrEmailDuplicate:
		writeError(w, http.StatusConflict, "auth.email_duplicate", err.Error())
	case err == models.ErrUserNotFound:
		writeError(w, http.StatusNotFound, "auth.user_not_found", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected error occurred")
	}
}

func mapErrorToCode(err error) string {
	switch err {
	case models.ErrNameRequired:
		return "auth.name_required"
	case models.ErrEmailRequired:
		return "auth.email_required"
	case models.ErrPasswordRequired:
		return "auth.password_required"
	case models.ErrPasswordTooShort:
		return "auth.password_too_short"
	default:
		return "auth.validation_error"
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
	json.NewEncoder(w).Encode(map[string]any{
		"error":   code,
		"message": message,
		"errors":  []string{code},
	})
}

func writeValidationErrors(w http.ResponseWriter, status int, code string, errors []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error":  code,
		"errors": errors,
	})
}
