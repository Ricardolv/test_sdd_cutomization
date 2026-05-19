package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/providers"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/repositories"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/services"
)

type UserHandler struct {
	repo   repositories.AuthRepository
	crypto providers.CryptoProvider
}

func NewUserHandler(repo repositories.AuthRepository, crypto providers.CryptoProvider) *UserHandler {
	return &UserHandler{repo: repo, crypto: crypto}
}

type saveUserRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req saveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	uc := services.NewSaveUser(h.repo, h.crypto)
	err := uc.Execute(r.Context(), services.SaveUserInput{
		ID:       req.ID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handleSaveUserError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req saveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}

	uc := services.NewSaveUser(h.repo, h.crypto)
	err := uc.Execute(r.Context(), services.SaveUserInput{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		handleSaveUserError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uc := services.NewDeleteUser(h.repo)
	err := uc.Execute(r.Context(), services.DeleteUserInput{ID: id})
	if err != nil {
		handleDeleteUserError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		handleDeleteUserError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	users, err := h.repo.ListPaginated(r.Context(), offset, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to list users")
		return
	}

	total, err := h.repo.Count(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to count users")
		return
	}

	items := make([]map[string]any, len(users))
	for i, u := range users {
		items[i] = map[string]any{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func handleSaveUserError(w http.ResponseWriter, err error) {
	switch {
	case err == models.ErrNameRequired:
		writeError(w, http.StatusUnprocessableEntity, "auth.name_required", "Name is required")
	case err == models.ErrEmailRequired:
		writeError(w, http.StatusUnprocessableEntity, "auth.email_required", "Email is required")
	case err == models.ErrPasswordRequired:
		writeError(w, http.StatusUnprocessableEntity, "auth.password_required", "Password is required")
	case err == models.ErrPasswordTooShort:
		writeError(w, http.StatusUnprocessableEntity, "auth.password_too_short", "Password must be at least 6 characters")
	case err == models.ErrEmailDuplicate:
		writeError(w, http.StatusConflict, "auth.email_duplicate", "This email is already registered")
	case err == models.ErrUserNotFound:
		writeError(w, http.StatusNotFound, "auth.user_not_found", "User not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected error occurred")
	}
}

func handleDeleteUserError(w http.ResponseWriter, err error) {
	switch {
	case err == models.ErrUserNotFound:
		writeError(w, http.StatusNotFound, "auth.user_not_found", "User not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "unexpected error occurred")
	}
}
