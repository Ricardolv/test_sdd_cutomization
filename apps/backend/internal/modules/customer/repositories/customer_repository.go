package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/models"
)

type CustomerRepository interface {
	Create(ctx context.Context, c *models.Customer) error
	Update(ctx context.Context, c *models.Customer) error
	FindByID(ctx context.Context, id string) (*models.Customer, error)
	FindByEmail(ctx context.Context, email string) (*models.Customer, error)
	List(ctx context.Context) ([]models.Customer, error)
	Deactivate(ctx context.Context, id string) error
}

type postgresCustomerRepository struct {
	db *sql.DB
}

func NewPostgresCustomerRepository(db *sql.DB) CustomerRepository {
	return &postgresCustomerRepository{db: db}
}

func (r *postgresCustomerRepository) Create(ctx context.Context, c *models.Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone, notes, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.Name, c.Email, c.Phone, c.Notes, c.Active, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return models.ErrEmailDuplicate
		}
		return err
	}
	return nil
}

func (r *postgresCustomerRepository) Update(ctx context.Context, c *models.Customer) error {
	query := `
		UPDATE customers
		SET name = $2, email = $3, phone = $4, notes = $5, active = $6, updated_at = $7
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		c.ID, c.Name, c.Email, c.Phone, c.Notes, c.Active, c.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return models.ErrEmailDuplicate
		}
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.ErrCustomerNotFound
	}

	return nil
}

func (r *postgresCustomerRepository) FindByID(ctx context.Context, id string) (*models.Customer, error) {
	query := `
		SELECT id, name, email, phone, notes, active, created_at, updated_at
		FROM customers
		WHERE id = $1 AND active = true
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanCustomer(row)
}

func (r *postgresCustomerRepository) FindByEmail(ctx context.Context, email string) (*models.Customer, error) {
	query := `
		SELECT id, name, email, phone, notes, active, created_at, updated_at
		FROM customers
		WHERE email = $1 AND active = true
	`
	row := r.db.QueryRowContext(ctx, query, email)
	return r.scanCustomer(row)
}

func (r *postgresCustomerRepository) List(ctx context.Context) ([]models.Customer, error) {
	query := `
		SELECT id, name, email, phone, notes, active, created_at, updated_at
		FROM customers
		WHERE active = true
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Notes, &c.Active, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *postgresCustomerRepository) Deactivate(ctx context.Context, id string) error {
	query := `
		UPDATE customers
		SET active = false, updated_at = $2
		WHERE id = $1 AND active = true
	`
	result, err := r.db.ExecContext(ctx, query, id, time.Now().UTC())
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.ErrCustomerNotFound
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func (r *postgresCustomerRepository) scanCustomer(row scanner) (*models.Customer, error) {
	var c models.Customer
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Notes, &c.Active, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrCustomerNotFound
		}
		return nil, err
	}
	return &c, nil
}
