package services

import (
	"context"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/catalog/repositories"
)

type SaveProductInput struct {
	ID              string
	Name            string
	Description     string
	Price           float64
	Status          models.ProductStatus
	AvailableOnline bool
	Featured        bool
	AllowsPreOrder  bool
}

type SaveProductUseCase struct {
	repo repositories.ProductRepository
}

func NewSaveProductUseCase(repo repositories.ProductRepository) *SaveProductUseCase {
	return &SaveProductUseCase{repo: repo}
}

func (u *SaveProductUseCase) Execute(ctx context.Context, input SaveProductInput) error {
	if input.ID != "" {
		existing, err := u.repo.FindByID(ctx, input.ID)
		if err == nil && existing != nil {
			return u.update(ctx, input.ID, input)
		}
	}

	return u.create(ctx, input)
}

func (u *SaveProductUseCase) create(ctx context.Context, input SaveProductInput) error {
	p, err := models.NewProduct(
		input.Name,
		input.Description,
		input.Price,
		input.Status,
		input.AvailableOnline,
		input.Featured,
		input.AllowsPreOrder,
	)
	if err != nil {
		return err
	}

	if input.ID != "" {
		p.ID = input.ID
	}

	return u.repo.Create(ctx, p)
}

func (u *SaveProductUseCase) update(ctx context.Context, id string, input SaveProductInput) error {
	p, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return models.ErrProductNotFound
	}

	if err := p.Update(
		input.Name,
		input.Description,
		input.Price,
		input.Status,
		input.AvailableOnline,
		input.Featured,
		input.AllowsPreOrder,
	); err != nil {
		return err
	}

	return u.repo.Update(ctx, p)
}
