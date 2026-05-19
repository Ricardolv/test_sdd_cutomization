package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
)

type AuthRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
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
