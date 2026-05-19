package services

import (
	"context"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/providers"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/repositories"
)

type AuthService struct {
	repo   repositories.AuthRepository
	crypto providers.CryptoProvider
}

func NewAuthService(repo repositories.AuthRepository, crypto providers.CryptoProvider) *AuthService {
	return &AuthService{repo: repo, crypto: crypto}
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

	hashedPassword, err := s.crypto.HashPassword(ctx, u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = hashedPassword

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}
