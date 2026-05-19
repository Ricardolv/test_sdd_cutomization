package services

import (
	"context"
	"testing"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/repositories"
)

type fakeCustomerRepo struct {
	customers map[string]*models.Customer
	byEmail   map[string]*models.Customer
}

func newFakeCustomerRepo() *fakeCustomerRepo {
	return &fakeCustomerRepo{
		customers: make(map[string]*models.Customer),
		byEmail:   make(map[string]*models.Customer),
	}
}

func (f *fakeCustomerRepo) Create(ctx context.Context, c *models.Customer) error {
	if _, exists := f.byEmail[c.Email]; exists {
		return models.ErrEmailDuplicate
	}
	f.customers[c.ID] = c
	f.byEmail[c.Email] = c
	return nil
}

func (f *fakeCustomerRepo) Update(ctx context.Context, c *models.Customer) error {
	if _, exists := f.customers[c.ID]; !exists {
		return models.ErrCustomerNotFound
	}
	delete(f.byEmail, c.Email)
	f.customers[c.ID] = c
	f.byEmail[c.Email] = c
	return nil
}

func (f *fakeCustomerRepo) FindByID(ctx context.Context, id string) (*models.Customer, error) {
	c, exists := f.customers[id]
	if !exists {
		return nil, models.ErrCustomerNotFound
	}
	return c, nil
}

func (f *fakeCustomerRepo) FindByEmail(ctx context.Context, email string) (*models.Customer, error) {
	c, exists := f.byEmail[email]
	if !exists {
		return nil, models.ErrCustomerNotFound
	}
	return c, nil
}

func (f *fakeCustomerRepo) List(ctx context.Context) ([]models.Customer, error) {
	var result []models.Customer
	for _, c := range f.customers {
		result = append(result, *c)
	}
	return result, nil
}

func (f *fakeCustomerRepo) Deactivate(ctx context.Context, id string) error {
	c, exists := f.customers[id]
	if !exists || !c.Active {
		return models.ErrCustomerNotFound
	}
	c.Deactivate()
	return nil
}

var _ repositories.CustomerRepository = (*fakeCustomerRepo)(nil)

func TestCustomerService_Save_Create(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	c, err := svc.Save(context.Background(), "", "John Doe", "john@example.com", "+1234567890", "VIP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Name != "John Doe" {
		t.Errorf("expected name John Doe, got %s", c.Name)
	}
	if c.Email != "john@example.com" {
		t.Errorf("expected email john@example.com, got %s", c.Email)
	}
	if !c.Active {
		t.Error("expected customer to be active")
	}
}

func TestCustomerService_Save_CreateMissingName(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	_, err := svc.Save(context.Background(), "", "", "john@example.com", "", "")
	if err != models.ErrNameRequired {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestCustomerService_Save_CreateMissingEmail(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	_, err := svc.Save(context.Background(), "", "John Doe", "", "", "")
	if err != models.ErrEmailRequired {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestCustomerService_Save_CreateDuplicateEmail(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	_, err := svc.Save(context.Background(), "", "John", "john@example.com", "", "")
	if err != nil {
		t.Fatalf("unexpected error on first create: %v", err)
	}

	_, err = svc.Save(context.Background(), "", "Jane", "john@example.com", "", "")
	if err != models.ErrEmailDuplicate {
		t.Fatalf("expected ErrEmailDuplicate, got %v", err)
	}
}

func TestCustomerService_Save_Update(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	created, _ := svc.Save(context.Background(), "", "John", "john@example.com", "", "")

	updated, err := svc.Save(context.Background(), created.ID, "John Updated", "john.updated@example.com", "+111", "Updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.Name != "John Updated" {
		t.Errorf("expected name John Updated, got %s", updated.Name)
	}
	if updated.Email != "john.updated@example.com" {
		t.Errorf("expected email john.updated@example.com, got %s", updated.Email)
	}
}

func TestCustomerService_Save_UpdateNotFound(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	_, err := svc.Save(context.Background(), "nonexistent-id", "John", "john@example.com", "", "")
	if err != models.ErrCustomerNotFound {
		t.Fatalf("expected ErrCustomerNotFound, got %v", err)
	}
}

func TestCustomerService_Save_UpdateDuplicateEmail(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	c1, _ := svc.Save(context.Background(), "", "John", "john@example.com", "", "")
	_, _ = svc.Save(context.Background(), "", "Jane", "jane@example.com", "", "")

	_, err := svc.Save(context.Background(), c1.ID, "John", "jane@example.com", "", "")
	if err != models.ErrEmailDuplicate {
		t.Fatalf("expected ErrEmailDuplicate, got %v", err)
	}
}

func TestCustomerService_Deactivate(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	created, _ := svc.Save(context.Background(), "", "John", "john@example.com", "", "")

	err := svc.Deactivate(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c, _ := repo.FindByID(context.Background(), created.ID)
	if c.Active {
		t.Error("expected customer to be deactivated")
	}
}

func TestCustomerService_DeactivateNotFound(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	err := svc.Deactivate(context.Background(), "nonexistent-id")
	if err != models.ErrCustomerNotFound {
		t.Fatalf("expected ErrCustomerNotFound, got %v", err)
	}
}

func TestCustomerService_GetByID(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	created, _ := svc.Save(context.Background(), "", "John", "john@example.com", "", "")

	c, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, c.ID)
	}
}

func TestCustomerService_GetByIDNotFound(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	_, err := svc.GetByID(context.Background(), "nonexistent-id")
	if err != models.ErrCustomerNotFound {
		t.Fatalf("expected ErrCustomerNotFound, got %v", err)
	}
}

func TestCustomerService_List(t *testing.T) {
	repo := newFakeCustomerRepo()
	svc := NewCustomerService(repo)

	_, _ = svc.Save(context.Background(), "", "John", "john@example.com", "", "")
	_, _ = svc.Save(context.Background(), "", "Jane", "jane@example.com", "", "")

	customers, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(customers) != 2 {
		t.Errorf("expected 2 customers, got %d", len(customers))
	}
}
