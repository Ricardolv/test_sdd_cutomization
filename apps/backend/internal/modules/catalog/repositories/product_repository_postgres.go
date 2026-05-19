package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/models"
)

type postgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(ctx context.Context, p *models.Product) error {
	query := `
		INSERT INTO products (id, name, description, price, status, available_online, featured, allows_pre_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.Name, p.Description, p.Price, p.Status, p.AvailableOnline, p.Featured, p.AllowsPreOrder, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *postgresProductRepository) Update(ctx context.Context, p *models.Product) error {
	query := `
		UPDATE products
		SET name = $2, description = $3, price = $4, status = $5, available_online = $6, featured = $7, allows_pre_order = $8, updated_at = $9
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		p.ID, p.Name, p.Description, p.Price, p.Status, p.AvailableOnline, p.Featured, p.AllowsPreOrder, p.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.ErrProductNotFound
	}

	return nil
}

func (r *postgresProductRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, status, available_online, featured, allows_pre_order, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanProduct(row)
}

func (r *postgresProductRepository) List(ctx context.Context, offset, limit int) ([]models.Product, int, error) {
	query := `
		SELECT id, name, description, price, status, available_online, featured, allows_pre_order, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Status, &p.AvailableOnline, &p.Featured, &p.AllowsPreOrder, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(*) FROM products`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *postgresProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.ErrProductNotFound
	}

	return nil
}

type productScanner interface {
	Scan(dest ...any) error
}

func (r *postgresProductRepository) scanProduct(row productScanner) (*models.Product, error) {
	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Status, &p.AvailableOnline, &p.Featured, &p.AllowsPreOrder, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}
