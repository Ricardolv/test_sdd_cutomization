package repositories

import (
	"context"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/models"
)

type ProductRepository interface {
	Create(ctx context.Context, p *models.Product) error
	Update(ctx context.Context, p *models.Product) error
	FindByID(ctx context.Context, id string) (*models.Product, error)
	List(ctx context.Context, offset, limit int) ([]models.Product, int, error)
	Delete(ctx context.Context, id string) error
}
