package services

import (
	"context"
	"testing"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/repositories"
)

type fakeProductRepo struct {
	products map[string]*models.Product
}

func newFakeProductRepo() *fakeProductRepo {
	return &fakeProductRepo{products: make(map[string]*models.Product)}
}

func (f *fakeProductRepo) Create(ctx context.Context, p *models.Product) error {
	f.products[p.ID] = p
	return nil
}

func (f *fakeProductRepo) Update(ctx context.Context, p *models.Product) error {
	if _, exists := f.products[p.ID]; !exists {
		return models.ErrProductNotFound
	}
	f.products[p.ID] = p
	return nil
}

func (f *fakeProductRepo) FindByID(ctx context.Context, id string) (*models.Product, error) {
	p, exists := f.products[id]
	if !exists {
		return nil, models.ErrProductNotFound
	}
	return p, nil
}

func (f *fakeProductRepo) List(ctx context.Context, offset, limit int) ([]models.Product, int, error) {
	var result []models.Product
	count := 0
	for _, p := range f.products {
		count++
		result = append(result, *p)
	}
	return result, count, nil
}

func (f *fakeProductRepo) Delete(ctx context.Context, id string) error {
	if _, exists := f.products[id]; !exists {
		return models.ErrProductNotFound
	}
	delete(f.products, id)
	return nil
}

var _ repositories.ProductRepository = (*fakeProductRepo)(nil)

func TestSaveProductUseCase_Create(t *testing.T) {
	repo := newFakeProductRepo()
	uc := NewSaveProductUseCase(repo)

	input := SaveProductInput{
		Name:            "Widget",
		Description:     "A useful widget",
		Price:           19.99,
		Status:          models.StatusActive,
		AvailableOnline: true,
		Featured:        false,
		AllowsPreOrder:  false,
	}

	err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(repo.products))
	}
}

func TestSaveProductUseCase_CreateWithID(t *testing.T) {
	repo := newFakeProductRepo()
	uc := NewSaveProductUseCase(repo)

	input := SaveProductInput{
		ID:              "custom-id-123",
		Name:            "Widget",
		Description:     "desc",
		Price:           10.00,
		Status:          models.StatusActive,
		AvailableOnline: false,
		Featured:        false,
		AllowsPreOrder:  false,
	}

	err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, exists := repo.products["custom-id-123"]
	if !exists {
		t.Fatal("expected product with custom ID to exist")
	}
	if p.Name != "Widget" {
		t.Errorf("expected name 'Widget', got %s", p.Name)
	}
}

func TestSaveProductUseCase_CreateValidationFailure(t *testing.T) {
	repo := newFakeProductRepo()
	uc := NewSaveProductUseCase(repo)

	input := SaveProductInput{
		Name:   "",
		Price:  10.00,
		Status: models.StatusActive,
	}

	err := uc.Execute(context.Background(), input)
	if err != models.ErrNameRequired {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestSaveProductUseCase_Update(t *testing.T) {
	repo := newFakeProductRepo()
	uc := NewSaveProductUseCase(repo)

	createInput := SaveProductInput{
		Name:   "Widget",
		Price:  10.00,
		Status: models.StatusActive,
	}
	_ = uc.Execute(context.Background(), createInput)

	var existingID string
	for id := range repo.products {
		existingID = id
		break
	}

	updateInput := SaveProductInput{
		ID:              existingID,
		Name:            "Gadget",
		Description:     "updated",
		Price:           25.50,
		Status:          models.StatusInactive,
		AvailableOnline: true,
		Featured:        true,
		AllowsPreOrder:  true,
	}

	err := uc.Execute(context.Background(), updateInput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := repo.products[existingID]
	if p.Name != "Gadget" {
		t.Errorf("expected name 'Gadget', got %s", p.Name)
	}
	if p.Price != 25.50 {
		t.Errorf("expected price 25.50, got %f", p.Price)
	}
	if p.Status != models.StatusInactive {
		t.Errorf("expected status inactive, got %s", p.Status)
	}
}

func TestSaveProductUseCase_UpdateNotFound(t *testing.T) {
	repo := newFakeProductRepo()
	uc := NewSaveProductUseCase(repo)

	input := SaveProductInput{
		ID:     "nonexistent-id",
		Name:   "Widget",
		Price:  10.00,
		Status: models.StatusActive,
	}

	err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.products) != 1 {
		t.Fatalf("expected 1 product (created), got %d", len(repo.products))
	}
}

func TestSaveProductUseCase_UpdateValidationFailure(t *testing.T) {
	repo := newFakeProductRepo()
	uc := NewSaveProductUseCase(repo)

	createInput := SaveProductInput{
		Name:   "Widget",
		Price:  10.00,
		Status: models.StatusActive,
	}
	_ = uc.Execute(context.Background(), createInput)

	var existingID string
	for id := range repo.products {
		existingID = id
		break
	}

	updateInput := SaveProductInput{
		ID:     existingID,
		Name:   "",
		Price:  10.00,
		Status: models.StatusActive,
	}

	err := uc.Execute(context.Background(), updateInput)
	if err != models.ErrNameRequired {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestDeleteProductUseCase_Success(t *testing.T) {
	repo := newFakeProductRepo()
	uc := NewDeleteProductUseCase(repo)

	p, _ := models.NewProduct("Widget", "desc", 10.00, models.StatusActive, false, false, false)
	_ = repo.Create(context.Background(), p)

	err := uc.Execute(context.Background(), p.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.products) != 0 {
		t.Fatalf("expected 0 products after delete, got %d", len(repo.products))
	}
}

func TestDeleteProductUseCase_NotFound(t *testing.T) {
	repo := newFakeProductRepo()
	uc := NewDeleteProductUseCase(repo)

	err := uc.Execute(context.Background(), "nonexistent-id")
	if err != models.ErrProductNotFound {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}
