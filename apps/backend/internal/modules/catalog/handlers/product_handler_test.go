package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/repositories"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/services"
)

type fakeHandlerRepo struct {
	products map[string]*models.Product
}

func newFakeHandlerRepo() *fakeHandlerRepo {
	return &fakeHandlerRepo{products: make(map[string]*models.Product)}
}

func (f *fakeHandlerRepo) Create(ctx context.Context, p *models.Product) error {
	f.products[p.ID] = p
	return nil
}

func (f *fakeHandlerRepo) Update(ctx context.Context, p *models.Product) error {
	if _, exists := f.products[p.ID]; !exists {
		return models.ErrProductNotFound
	}
	f.products[p.ID] = p
	return nil
}

func (f *fakeHandlerRepo) FindByID(ctx context.Context, id string) (*models.Product, error) {
	p, exists := f.products[id]
	if !exists {
		return nil, models.ErrProductNotFound
	}
	return p, nil
}

func (f *fakeHandlerRepo) List(ctx context.Context, offset, limit int) ([]models.Product, int, error) {
	var result []models.Product
	for _, p := range f.products {
		result = append(result, *p)
	}
	return result, len(result), nil
}

func (f *fakeHandlerRepo) Delete(ctx context.Context, id string) error {
	if _, exists := f.products[id]; !exists {
		return models.ErrProductNotFound
	}
	delete(f.products, id)
	return nil
}

var _ repositories.ProductRepository = (*fakeHandlerRepo)(nil)

func setupProductHandler(t *testing.T) (*ProductHandler, *fakeHandlerRepo) {
	t.Helper()
	repo := newFakeHandlerRepo()
	saveUC := services.NewSaveProductUseCase(repo)
	deleteUC := services.NewDeleteProductUseCase(repo)
	handler := NewProductHandler(saveUC, deleteUC, repo)
	return handler, repo
}

func productBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(v)
	return &buf
}

func decodeProductResponse(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()
	json.NewDecoder(w.Body).Decode(v)
}

func TestProductHandler_Create_Success(t *testing.T) {
	handler, _ := setupProductHandler(t)

	req := CreateProductRequest{
		Name:            "Widget",
		Description:     "A useful widget",
		Price:           19.99,
		Status:          "active",
		AvailableOnline: true,
		Featured:        false,
		AllowsPreOrder:  false,
	}

	r := httptest.NewRequest(http.MethodPost, "/products", productBody(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestProductHandler_Create_MissingName(t *testing.T) {
	handler, _ := setupProductHandler(t)

	req := CreateProductRequest{
		Name:   "",
		Price:  10.00,
		Status: "active",
	}

	r := httptest.NewRequest(http.MethodPost, "/products", productBody(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_Create_NegativePrice(t *testing.T) {
	handler, _ := setupProductHandler(t)

	req := CreateProductRequest{
		Name:   "Widget",
		Price:  -5.00,
		Status: "active",
	}

	r := httptest.NewRequest(http.MethodPost, "/products", productBody(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_Create_InvalidStatus(t *testing.T) {
	handler, _ := setupProductHandler(t)

	req := CreateProductRequest{
		Name:   "Widget",
		Price:  10.00,
		Status: "unknown",
	}

	r := httptest.NewRequest(http.MethodPost, "/products", productBody(t, req))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_Create_InvalidJSON(t *testing.T) {
	handler, _ := setupProductHandler(t)

	r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.Create(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_Update_Success(t *testing.T) {
	handler, repo := setupProductHandler(t)
	router := handler.Routes()

	createReq := CreateProductRequest{
		Name:   "Widget",
		Price:  10.00,
		Status: "active",
	}
	r := httptest.NewRequest(http.MethodPost, "/", productBody(t, createReq))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	var createdID string
	for id := range repo.products {
		createdID = id
		break
	}

	updateReq := UpdateProductRequest{
		Name:            "Gadget",
		Description:     "updated",
		Price:           25.50,
		Status:          "inactive",
		AvailableOnline: true,
		Featured:        true,
		AllowsPreOrder:  false,
	}

	r = httptest.NewRequest(http.MethodPut, "/"+createdID, productBody(t, updateReq))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestProductHandler_Update_NotFound_CreatesNew(t *testing.T) {
	handler, repo := setupProductHandler(t)
	router := handler.Routes()

	req := UpdateProductRequest{
		Name:   "Widget",
		Price:  10.00,
		Status: "active",
	}

	r := httptest.NewRequest(http.MethodPut, "/nonexistent-id", productBody(t, req))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 (creates new with given ID), got %d", w.Code)
	}

	if len(repo.products) != 1 {
		t.Fatalf("expected 1 product created, got %d", len(repo.products))
	}
}

func TestProductHandler_GetByID_Success(t *testing.T) {
	handler, repo := setupProductHandler(t)
	router := handler.Routes()

	createReq := CreateProductRequest{
		Name:            "Widget",
		Description:     "desc",
		Price:           10.00,
		Status:          "active",
		AvailableOnline: true,
		Featured:        false,
		AllowsPreOrder:  false,
	}
	r := httptest.NewRequest(http.MethodPost, "/", productBody(t, createReq))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	var createdID string
	for id := range repo.products {
		createdID = id
		break
	}

	r = httptest.NewRequest(http.MethodGet, "/"+createdID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ProductResponse
	decodeProductResponse(t, w, &resp)

	if resp.ID != createdID {
		t.Errorf("expected ID %s, got %s", createdID, resp.ID)
	}
	if resp.Name != "Widget" {
		t.Errorf("expected name 'Widget', got %s", resp.Name)
	}
	if resp.Price != 10.00 {
		t.Errorf("expected price 10.00, got %f", resp.Price)
	}
	if resp.Status != "active" {
		t.Errorf("expected status 'active', got %s", resp.Status)
	}
	if !resp.AvailableOnline {
		t.Error("expected AvailableOnline to be true")
	}
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	handler, _ := setupProductHandler(t)

	r := httptest.NewRequest(http.MethodGet, "/products/nonexistent-id", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestProductHandler_List_Success(t *testing.T) {
	handler, repo := setupProductHandler(t)

	p, _ := models.NewProduct("Widget", "desc", 10.00, models.StatusActive, false, false, false)
	_ = repo.Create(context.Background(), p)
	p2, _ := models.NewProduct("Gadget", "desc", 20.00, models.StatusDraft, false, false, false)
	_ = repo.Create(context.Background(), p2)

	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.List(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListResponse
	decodeProductResponse(t, w, &resp)

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestProductHandler_List_Empty(t *testing.T) {
	handler, _ := setupProductHandler(t)

	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.List(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp ListResponse
	decodeProductResponse(t, w, &resp)

	if resp.Total != 0 {
		t.Errorf("expected total 0, got %d", resp.Total)
	}
}

func TestProductHandler_Delete_Success(t *testing.T) {
	handler, repo := setupProductHandler(t)
	router := handler.Routes()

	createReq := CreateProductRequest{
		Name:   "Widget",
		Price:  10.00,
		Status: "active",
	}
	r := httptest.NewRequest(http.MethodPost, "/", productBody(t, createReq))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	var createdID string
	for id := range repo.products {
		createdID = id
		break
	}

	r = httptest.NewRequest(http.MethodDelete, "/"+createdID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}

	if len(repo.products) != 0 {
		t.Fatalf("expected 0 products after delete, got %d", len(repo.products))
	}
}

func TestProductHandler_Delete_NotFound(t *testing.T) {
	handler, _ := setupProductHandler(t)

	r := httptest.NewRequest(http.MethodDelete, "/products/nonexistent-id", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestProductHandler_Routes(t *testing.T) {
	handler, _ := setupProductHandler(t)
	router := handler.Routes()

	if router == nil {
		t.Fatal("expected router to not be nil")
	}

	r := httptest.NewRequest(http.MethodPost, "/", productBody(t, CreateProductRequest{
		Name:   "Widget",
		Price:  10.00,
		Status: "active",
	}))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestToResponse(t *testing.T) {
	p := &models.Product{
		ID:              "test-id",
		Name:            "Widget",
		Description:     "desc",
		Price:           19.99,
		Status:          models.StatusActive,
		AvailableOnline: true,
		Featured:        false,
		AllowsPreOrder:  true,
	}

	resp := toResponse(p)

	if resp.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got %s", resp.ID)
	}
	if resp.Name != "Widget" {
		t.Errorf("expected name 'Widget', got %s", resp.Name)
	}
	if resp.Price != 19.99 {
		t.Errorf("expected price 19.99, got %f", resp.Price)
	}
	if resp.Status != "active" {
		t.Errorf("expected status 'active', got %s", resp.Status)
	}
	if !resp.AvailableOnline {
		t.Error("expected AvailableOnline to be true")
	}
	if resp.Featured {
		t.Error("expected Featured to be false")
	}
	if !resp.AllowsPreOrder {
		t.Error("expected AllowsPreOrder to be true")
	}
}

func TestHandleServiceError_ProductNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	handleServiceError(w, models.ErrProductNotFound)

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

func TestHandleServiceError_NameTooShort(t *testing.T) {
	w := httptest.NewRecorder()
	handleServiceError(w, models.ErrNameTooShort)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleServiceError_PricePrecision(t *testing.T) {
	w := httptest.NewRecorder()
	handleServiceError(w, models.ErrPricePrecision)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleServiceError_StatusInvalid(t *testing.T) {
	w := httptest.NewRecorder()
	handleServiceError(w, models.ErrStatusInvalid)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
