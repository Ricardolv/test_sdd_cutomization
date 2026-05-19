package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/repositories"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/services"
)

type fakeHandlerRepo struct {
	customers map[string]*models.Customer
	byEmail   map[string]*models.Customer
}

func newFakeHandlerRepo() *fakeHandlerRepo {
	return &fakeHandlerRepo{
		customers: make(map[string]*models.Customer),
		byEmail:   make(map[string]*models.Customer),
	}
}

func (f *fakeHandlerRepo) Create(ctx context.Context, c *models.Customer) error {
	if _, exists := f.byEmail[c.Email]; exists {
		return models.ErrEmailDuplicate
	}
	f.customers[c.ID] = c
	f.byEmail[c.Email] = c
	return nil
}

func (f *fakeHandlerRepo) Update(ctx context.Context, c *models.Customer) error {
	if _, exists := f.customers[c.ID]; !exists {
		return models.ErrCustomerNotFound
	}
	delete(f.byEmail, c.Email)
	f.customers[c.ID] = c
	f.byEmail[c.Email] = c
	return nil
}

func (f *fakeHandlerRepo) FindByID(ctx context.Context, id string) (*models.Customer, error) {
	c, exists := f.customers[id]
	if !exists || !c.Active {
		return nil, models.ErrCustomerNotFound
	}
	return c, nil
}

func (f *fakeHandlerRepo) FindByEmail(ctx context.Context, email string) (*models.Customer, error) {
	c, exists := f.byEmail[email]
	if !exists {
		return nil, models.ErrCustomerNotFound
	}
	return c, nil
}

func (f *fakeHandlerRepo) List(ctx context.Context) ([]models.Customer, error) {
	var result []models.Customer
	for _, c := range f.customers {
		result = append(result, *c)
	}
	return result, nil
}

func (f *fakeHandlerRepo) Deactivate(ctx context.Context, id string) error {
	c, exists := f.customers[id]
	if !exists || !c.Active {
		return models.ErrCustomerNotFound
	}
	c.Deactivate()
	return nil
}

var _ repositories.CustomerRepository = (*fakeHandlerRepo)(nil)

func setupHandler(t *testing.T) (*CustomerHandler, *fakeHandlerRepo) {
	t.Helper()
	repo := newFakeHandlerRepo()
	svc := services.NewCustomerService(repo)
	handler := NewCustomerHandler(svc)
	return handler, repo
}

func body(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(v)
	return &buf
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()
	json.NewDecoder(w.Body).Decode(v)
}

func TestCustomerHandler_Create_Success(t *testing.T) {
	handler, _ := setupHandler(t)

	req := CreateCustomerRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "+1234567890",
		Notes: "VIP",
	}

	r := httptest.NewRequest(http.MethodPost, "/customers", body(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp CustomerResponse
	decodeResponse(t, w, &resp)

	if resp.Name != "John Doe" {
		t.Errorf("expected name 'John Doe', got %s", resp.Name)
	}
	if resp.Email != "john@example.com" {
		t.Errorf("expected email 'john@example.com', got %s", resp.Email)
	}
	if !resp.Active {
		t.Error("expected customer to be active")
	}
}

func TestCustomerHandler_Create_MissingName(t *testing.T) {
	handler, _ := setupHandler(t)

	req := CreateCustomerRequest{
		Name:  "",
		Email: "john@example.com",
	}

	r := httptest.NewRequest(http.MethodPost, "/customers", body(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Create_MissingEmail(t *testing.T) {
	handler, _ := setupHandler(t)

	req := CreateCustomerRequest{
		Name:  "John Doe",
		Email: "",
	}

	r := httptest.NewRequest(http.MethodPost, "/customers", body(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Create_DuplicateEmail(t *testing.T) {
	handler, _ := setupHandler(t)

	req1 := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r1 := httptest.NewRequest(http.MethodPost, "/customers", body(t, req1))
	w1 := httptest.NewRecorder()
	handler.Create(w1, r1)

	req2 := CreateCustomerRequest{
		Name:  "Jane",
		Email: "john@example.com",
	}
	r2 := httptest.NewRequest(http.MethodPost, "/customers", body(t, req2))
	w2 := httptest.NewRecorder()
	handler.Create(w2, r2)

	if w2.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", w2.Code)
	}
}

func TestCustomerHandler_Create_InvalidJSON(t *testing.T) {
	handler, _ := setupHandler(t)

	r := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Update_Success(t *testing.T) {
	handler, _ := setupHandler(t)
	router := handler.Routes()

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	var created CustomerResponse
	decodeResponse(t, w, &created)

	updateReq := UpdateCustomerRequest{
		Name:  "John Updated",
		Email: "john.updated@example.com",
		Phone: "+999",
		Notes: "Updated notes",
	}

	r = httptest.NewRequest(http.MethodPut, "/"+created.ID, body(t, updateReq))
	w = httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp CustomerResponse
	decodeResponse(t, w, &resp)

	if resp.Name != "John Updated" {
		t.Errorf("expected name 'John Updated', got %s", resp.Name)
	}
	if resp.Email != "john.updated@example.com" {
		t.Errorf("expected email 'john.updated@example.com', got %s", resp.Email)
	}
}

func TestCustomerHandler_Update_NotFound(t *testing.T) {
	handler, _ := setupHandler(t)
	router := handler.Routes()

	req := UpdateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}

	r := httptest.NewRequest(http.MethodPut, "/customers/nonexistent-id", body(t, req))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestCustomerHandler_Update_DuplicateEmail(t *testing.T) {
	handler, _ := setupHandler(t)
	router := handler.Routes()

	req1 := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req1))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	var c1 CustomerResponse
	decodeResponse(t, w, &c1)

	req2 := CreateCustomerRequest{
		Name:  "Jane",
		Email: "jane@example.com",
	}
	r = httptest.NewRequest(http.MethodPost, "/", body(t, req2))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	updateReq := UpdateCustomerRequest{
		Name:  "John",
		Email: "jane@example.com",
	}

	r = httptest.NewRequest(http.MethodPut, "/"+c1.ID, body(t, updateReq))
	w = httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", w.Code)
	}
}

func TestCustomerHandler_Update_MissingName(t *testing.T) {
	handler, _ := setupHandler(t)
	router := handler.Routes()

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	var created CustomerResponse
	decodeResponse(t, w, &created)

	updateReq := UpdateCustomerRequest{
		Name:  "",
		Email: "john@example.com",
	}

	r = httptest.NewRequest(http.MethodPut, "/"+created.ID, body(t, updateReq))
	w = httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCustomerHandler_GetByID_Success(t *testing.T) {
	handler, _ := setupHandler(t)
	router := handler.Routes()

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
		Phone: "+123",
		Notes: "notes",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	var created CustomerResponse
	decodeResponse(t, w, &created)

	r = httptest.NewRequest(http.MethodGet, "/"+created.ID, nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp CustomerResponse
	decodeResponse(t, w, &resp)

	if resp.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, resp.ID)
	}
	if resp.Name != "John" {
		t.Errorf("expected name 'John', got %s", resp.Name)
	}
}

func TestCustomerHandler_GetByID_NotFound(t *testing.T) {
	handler, _ := setupHandler(t)

	r := httptest.NewRequest(http.MethodGet, "/customers/nonexistent-id", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestCustomerHandler_List_Success(t *testing.T) {
	handler, repo := setupHandler(t)
	svc := services.NewCustomerService(repo)

	_, _ = svc.Save(context.Background(), "", "John", "john@example.com", "", "")
	_, _ = svc.Save(context.Background(), "", "Jane", "jane@example.com", "", "")

	r := httptest.NewRequest(http.MethodGet, "/customers", nil)
	w := httptest.NewRecorder()

	handler.List(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp []CustomerResponse
	decodeResponse(t, w, &resp)

	if len(resp) != 2 {
		t.Errorf("expected 2 customers, got %d", len(resp))
	}
}

func TestCustomerHandler_List_Empty(t *testing.T) {
	handler, _ := setupHandler(t)

	r := httptest.NewRequest(http.MethodGet, "/customers", nil)
	w := httptest.NewRecorder()

	handler.List(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp []CustomerResponse
	decodeResponse(t, w, &resp)

	if len(resp) != 0 {
		t.Errorf("expected 0 customers, got %d", len(resp))
	}
}

func TestCustomerHandler_Deactivate_Success(t *testing.T) {
	handler, _ := setupHandler(t)
	router := handler.Routes()

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	var created CustomerResponse
	decodeResponse(t, w, &created)

	r = httptest.NewRequest(http.MethodDelete, "/"+created.ID, nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}

	r = httptest.NewRequest(http.MethodGet, "/"+created.ID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected customer to be deactivated (not found), got %d", w.Code)
	}
}

func TestCustomerHandler_Deactivate_NotFound(t *testing.T) {
	handler, _ := setupHandler(t)

	r := httptest.NewRequest(http.MethodDelete, "/customers/nonexistent-id", nil)
	w := httptest.NewRecorder()

	handler.Deactivate(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestCustomerHandler_Routes(t *testing.T) {
	handler, _ := setupHandler(t)
	router := handler.Routes()

	if router == nil {
		t.Fatal("expected router to not be nil")
	}

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":"John","email":"john@example.com"}`)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestCustomerHandler_Create_ResponseHasAllFields(t *testing.T) {
	handler, _ := setupHandler(t)

	req := CreateCustomerRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "+123",
		Notes: "Test notes",
	}

	r := httptest.NewRequest(http.MethodPost, "/customers", body(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	var resp CustomerResponse
	decodeResponse(t, w, &resp)

	if resp.ID == "" {
		t.Error("expected ID to be set")
	}
	if resp.Phone != "+123" {
		t.Errorf("expected phone '+123', got %s", resp.Phone)
	}
	if resp.Notes != "Test notes" {
		t.Errorf("expected notes 'Test notes', got %s", resp.Notes)
	}
	if resp.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if resp.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestCustomerHandler_Update_ResponseHasUpdatedTimestamp(t *testing.T) {
	handler, _ := setupHandler(t)
	router := handler.Routes()

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	var created CustomerResponse
	decodeResponse(t, w, &created)

	time.Sleep(10 * time.Millisecond)

	updateReq := UpdateCustomerRequest{
		Name:  "John Updated",
		Email: "john@example.com",
	}

	r = httptest.NewRequest(http.MethodPut, "/"+created.ID, body(t, updateReq))
	w = httptest.NewRecorder()

	router.ServeHTTP(w, r)

	var resp CustomerResponse
	decodeResponse(t, w, &resp)

	if !resp.UpdatedAt.After(created.CreatedAt) {
		t.Errorf("expected UpdatedAt %v to be after CreatedAt %v", resp.UpdatedAt, created.CreatedAt)
	}
}

func TestCustomerHandler_ErrorResponseFormat(t *testing.T) {
	handler, _ := setupHandler(t)

	req := CreateCustomerRequest{
		Name:  "",
		Email: "john@example.com",
	}

	r := httptest.NewRequest(http.MethodPost, "/customers", body(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	var errResp map[string]string
	decodeResponse(t, w, &errResp)

	if _, ok := errResp["error"]; !ok {
		t.Error("expected error response to have 'error' field")
	}
	if _, ok := errResp["message"]; !ok {
		t.Error("expected error response to have 'message' field")
	}
}

func TestCustomerHandler_Create_ContentType(t *testing.T) {
	handler, _ := setupHandler(t)

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}

	r := httptest.NewRequest(http.MethodPost, "/customers", body(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %s", w.Header().Get("Content-Type"))
	}
}

func TestCustomerHandler_List_ContentType(t *testing.T) {
	handler, _ := setupHandler(t)

	r := httptest.NewRequest(http.MethodGet, "/customers", nil)
	w := httptest.NewRecorder()

	handler.List(w, r)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %s", w.Header().Get("Content-Type"))
	}
}

func TestCustomerHandler_Update_InvalidJSON(t *testing.T) {
	handler, _ := setupHandler(t)

	r := httptest.NewRequest(http.MethodPut, "/customers/some-id", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	handler.Update(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCustomerHandler_Routes_AllMethods(t *testing.T) {
	handler, _ := setupHandler(t)
	router := handler.Routes()

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	var created CustomerResponse
	decodeResponse(t, w, &created)

	tests := []struct {
		method string
		path   string
		body   *bytes.Buffer
		expect int
	}{
		{http.MethodPost, "/", body(t, map[string]string{"name": "Jane", "email": "jane@example.com"}), http.StatusCreated},
		{http.MethodGet, "/", nil, http.StatusOK},
		{http.MethodGet, "/" + created.ID, nil, http.StatusOK},
		{http.MethodPut, "/" + created.ID, body(t, map[string]string{"name": "Updated", "email": "john@example.com"}), http.StatusOK},
		{http.MethodDelete, "/" + created.ID, nil, http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			var r *http.Request
			if tt.body != nil {
				r = httptest.NewRequest(tt.method, tt.path, tt.body)
			} else {
				r = httptest.NewRequest(tt.method, tt.path, nil)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			if w.Code != tt.expect {
				t.Errorf("expected status %d, got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestToResponse(t *testing.T) {
	now := time.Now().UTC()
	c := &models.Customer{
		ID:        "test-id",
		Name:      "Test",
		Email:     "test@example.com",
		Phone:     "+123",
		Notes:     "notes",
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := toResponse(c)

	if resp.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got %s", resp.ID)
	}
	if resp.Name != "Test" {
		t.Errorf("expected name 'Test', got %s", resp.Name)
	}
	if resp.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", resp.Email)
	}
	if resp.Phone != "+123" {
		t.Errorf("expected phone '+123', got %s", resp.Phone)
	}
	if resp.Notes != "notes" {
		t.Errorf("expected notes 'notes', got %s", resp.Notes)
	}
	if !resp.Active {
		t.Error("expected Active to be true")
	}
	if !resp.CreatedAt.Equal(now) {
		t.Error("expected CreatedAt to match")
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %s", w.Header().Get("Content-Type"))
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["key"] != "value" {
		t.Errorf("expected 'value', got %s", resp["key"])
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusBadRequest, "test_error", "test message")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["error"] != "test_error" {
		t.Errorf("expected error 'test_error', got %s", resp["error"])
	}
	if resp["message"] != "test message" {
		t.Errorf("expected message 'test message', got %s", resp["message"])
	}
}

func TestHandleServiceError_NameRequired(t *testing.T) {
	w := httptest.NewRecorder()
	handleServiceError(w, models.ErrNameRequired)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleServiceError_EmailRequired(t *testing.T) {
	w := httptest.NewRecorder()
	handleServiceError(w, models.ErrEmailRequired)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleServiceError_EmailDuplicate(t *testing.T) {
	w := httptest.NewRecorder()
	handleServiceError(w, models.ErrEmailDuplicate)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", w.Code)
	}
}

func TestHandleServiceError_CustomerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	handleServiceError(w, models.ErrCustomerNotFound)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestHandleServiceError_UnknownError(t *testing.T) {
	w := httptest.NewRecorder()
	handleServiceError(w, context.DeadlineExceeded)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestCustomerHandler_GetByID_WithChiRouter(t *testing.T) {
	handler, _ := setupHandler(t)

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req))
	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, r)
	var created CustomerResponse
	decodeResponse(t, w, &created)

	router := chi.NewRouter()
	router.Get("/{id}", handler.GetByID)

	req2 := httptest.NewRequest(http.MethodGet, "/"+created.ID, nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req2)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCustomerHandler_Update_WithChiRouter(t *testing.T) {
	handler, _ := setupHandler(t)

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req))
	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, r)
	var created CustomerResponse
	decodeResponse(t, w, &created)

	router := chi.NewRouter()
	router.Put("/{id}", handler.Update)

	req2 := httptest.NewRequest(http.MethodPut, "/"+created.ID, body(t, UpdateCustomerRequest{
		Name:  "Updated",
		Email: "john@example.com",
	}))
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req2)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCustomerHandler_Deactivate_WithChiRouter(t *testing.T) {
	handler, _ := setupHandler(t)

	req := CreateCustomerRequest{
		Name:  "John",
		Email: "john@example.com",
	}
	r := httptest.NewRequest(http.MethodPost, "/", body(t, req))
	w := httptest.NewRecorder()
	handler.Routes().ServeHTTP(w, r)
	var created CustomerResponse
	decodeResponse(t, w, &created)

	router := chi.NewRouter()
	router.Delete("/{id}", handler.Deactivate)

	req2 := httptest.NewRequest(http.MethodDelete, "/"+created.ID, nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req2)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}
