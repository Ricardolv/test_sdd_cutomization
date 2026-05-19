package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
)

type AuthRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id string) error
	ListPaginated(ctx context.Context, offset, limit int) ([]models.User, error)
	Count(ctx context.Context) (int, error)
}

type postgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) AuthRepository {
	return &postgresAuthRepository{db: db}
}

func (r *postgresAuthRepository) Create(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (id, name, email, password, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Name, u.Email, u.Password, u.Active, u.CreatedAt, u.UpdatedAt,
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

func (r *postgresAuthRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password, active, created_at, updated_at
		FROM users
		WHERE email = $1 AND active = true
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresAuthRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, name, email, password, active, created_at, updated_at
		FROM users
		WHERE id = $1 AND active = true
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresAuthRepository) Update(ctx context.Context, u *models.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, password = $4, active = $5, updated_at = $6
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		u.ID, u.Name, u.Email, u.Password, u.Active, u.UpdatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return models.ErrEmailDuplicate
		}
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrUserNotFound
	}
	return nil
}

func (r *postgresAuthRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET active = false, updated_at = $2
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id, time.Now().UTC())
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrUserNotFound
	}
	return nil
}

func (r *postgresAuthRepository) ListPaginated(ctx context.Context, offset, limit int) ([]models.User, error) {
	query := `
		SELECT id, name, email, active, created_at, updated_at
		FROM users
		WHERE active = true
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Active, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *postgresAuthRepository) Count(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE active = true
	`
	var count int
	if err := r.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
