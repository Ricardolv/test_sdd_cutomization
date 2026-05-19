package services

import (
	"context"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/repositories"
)

type DeleteProductUseCase struct {
	repo repositories.ProductRepository
}

func NewDeleteProductUseCase(repo repositories.ProductRepository) *DeleteProductUseCase {
	return &DeleteProductUseCase{repo: repo}
}

func (u *DeleteProductUseCase) Execute(ctx context.Context, id string) error {
	_, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return models.ErrProductNotFound
	}

	return u.repo.Delete(ctx, id)
}
