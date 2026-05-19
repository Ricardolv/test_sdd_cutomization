package services

import (
	"context"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/repositories"
)

type AuthService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	u, err := models.NewUser(name, email, password)
	if err != nil {
		return nil, err
	}

	existing, _ := s.repo.FindByEmail(ctx, u.Email)
	if existing != nil {
		return nil, models.ErrEmailDuplicate
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}
