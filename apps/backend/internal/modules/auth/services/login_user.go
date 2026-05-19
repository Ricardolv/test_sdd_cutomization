package services

import (
	"context"
	"errors"
	"regexp"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/providers"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/repositories"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	ID    string
	Name  string
	Email string
}

type LoginUser struct {
	repo   repositories.AuthRepository
	crypto providers.CryptoProvider
}

func NewLoginUser(repo repositories.AuthRepository, crypto providers.CryptoProvider) *LoginUser {
	return &LoginUser{repo: repo, crypto: crypto}
}

func (uc *LoginUser) Execute(ctx context.Context, input LoginUserInput) (*LoginUserOutput, error) {
	if input.Email == "" {
		return nil, errors.New("email is required")
	}
	if !emailRegex.MatchString(input.Email) {
		return nil, errors.New("email is invalid")
	}
	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}

	if err := uc.crypto.ComparePassword(ctx, input.Password, user.Password); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	return &LoginUserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
