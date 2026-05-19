package services

import (
	"context"
	"errors"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/customer/repositories"
)

type CustomerService struct {
	repo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Save(ctx context.Context, id, name, email, phone, notes string) (*models.Customer, error) {
	if id != "" {
		return s.update(ctx, id, name, email, phone, notes)
	}
	return s.create(ctx, name, email, phone, notes)
}

func (s *CustomerService) create(ctx context.Context, name, email, phone, notes string) (*models.Customer, error) {
	c, err := models.NewCustomer(name, email, phone, notes)
	if err != nil {
		return nil, err
	}

	existing, _ := s.repo.FindByEmail(ctx, c.Email)
	if existing != nil {
		return nil, models.ErrEmailDuplicate
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *CustomerService) update(ctx context.Context, id, name, email, phone, notes string) (*models.Customer, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, models.ErrCustomerNotFound
	}

	if err := c.Update(name, email, phone, notes); err != nil {
		return nil, err
	}

	existing, _ := s.repo.FindByEmail(ctx, c.Email)
	if existing != nil && existing.ID != c.ID {
		return nil, models.ErrEmailDuplicate
	}

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *CustomerService) Deactivate(ctx context.Context, id string) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return models.ErrCustomerNotFound
	}

	c.Deactivate()
	return s.repo.Update(ctx, c)
}

func (s *CustomerService) GetByID(ctx context.Context, id string) (*models.Customer, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, models.ErrCustomerNotFound
	}
	return c, nil
}

func (s *CustomerService) List(ctx context.Context) ([]models.Customer, error) {
	return s.repo.List(ctx)
}

func IsDuplicateEmail(err error) bool {
	return errors.Is(err, models.ErrEmailDuplicate)
}
