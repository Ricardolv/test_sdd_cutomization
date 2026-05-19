package repositories

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/google/uuid"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/models"
)

func testDB(t *testing.T) *postgresCustomerRepository {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/testapp?sslmode=disable"
	}

	t.Helper()

	db, err := setupTestDB(databaseURL)
	if err != nil {
		t.Skipf("database not available: %v", err)
	}

	t.Cleanup(func() {
		db.Exec("DELETE FROM customers")
	})

	return NewPostgresCustomerRepository(db).(*postgresCustomerRepository)
}

func setupTestDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			phone VARCHAR(50),
			notes TEXT,
			active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`)

	return db, nil
}

func newTestCustomer() *models.Customer {
	return &models.Customer{
		ID:        uuid.New().String(),
		Name:      "Test Customer",
		Email:     "test-" + uuid.New().String()[:8] + "@example.com",
		Phone:     "+1234567890",
		Notes:     "Test notes",
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func TestPostgresCustomerRepository_Create(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()
	c := newTestCustomer()

	err := repo.Create(ctx, c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByID(ctx, c.ID)
	if err != nil {
		t.Fatalf("unexpected error finding customer: %v", err)
	}

	if found.Name != c.Name {
		t.Errorf("expected name %s, got %s", c.Name, found.Name)
	}
}

func TestPostgresCustomerRepository_CreateDuplicateEmail(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()
	c1 := newTestCustomer()
	c2 := newTestCustomer()
	c2.Email = c1.Email

	_ = repo.Create(ctx, c1)
	err := repo.Create(ctx, c2)

	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
	if !errors.Is(err, models.ErrEmailDuplicate) {
		t.Fatalf("expected ErrEmailDuplicate, got %v", err)
	}
}

func TestPostgresCustomerRepository_Update(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()
	c := newTestCustomer()

	_ = repo.Create(ctx, c)

	c.Name = "Updated Name"
	c.UpdatedAt = time.Now().UTC()

	err := repo.Update(ctx, c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := repo.FindByID(ctx, c.ID)
	if found.Name != "Updated Name" {
		t.Errorf("expected name Updated Name, got %s", found.Name)
	}
}

func TestPostgresCustomerRepository_UpdateNotFound(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()
	c := newTestCustomer()
	c.ID = uuid.New().String()

	err := repo.Update(ctx, c)
	if err == nil {
		t.Fatal("expected error for non-existent customer")
	}
}

func TestPostgresCustomerRepository_FindByID(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()
	c := newTestCustomer()

	_ = repo.Create(ctx, c)

	found, err := repo.FindByID(ctx, c.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email != c.Email {
		t.Errorf("expected email %s, got %s", c.Email, found.Email)
	}
}

func TestPostgresCustomerRepository_FindByIDNotFound(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent customer")
	}
	if !errors.Is(err, models.ErrCustomerNotFound) {
		t.Fatalf("expected ErrCustomerNotFound, got %v", err)
	}
}

func TestPostgresCustomerRepository_FindByEmail(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()
	c := newTestCustomer()

	_ = repo.Create(ctx, c)

	found, err := repo.FindByEmail(ctx, c.Email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.ID != c.ID {
		t.Errorf("expected ID %s, got %s", c.ID, found.ID)
	}
}

func TestPostgresCustomerRepository_List(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()

	c1 := newTestCustomer()
	c2 := newTestCustomer()

	_ = repo.Create(ctx, c1)
	_ = repo.Create(ctx, c2)

	customers, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(customers) != 2 {
		t.Errorf("expected 2 customers, got %d", len(customers))
	}
}

func TestPostgresCustomerRepository_Deactivate(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()
	c := newTestCustomer()

	_ = repo.Create(ctx, c)

	err := repo.Deactivate(ctx, c.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.FindByID(ctx, c.ID)
	if err == nil {
		t.Fatal("expected error finding deactivated customer")
	}
}

func TestPostgresCustomerRepository_DeactivateNotFound(t *testing.T) {
	repo := testDB(t)
	ctx := context.Background()

	err := repo.Deactivate(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent customer")
	}
}
