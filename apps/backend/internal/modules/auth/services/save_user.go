package services

import (
	"context"
	"strings"
	"time"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/providers"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/repositories"
)

type SaveUserInput struct {
	ID       string
	Name     string
	Email    string
	Password string
}

type SaveUser struct {
	repo   repositories.AuthRepository
	crypto providers.CryptoProvider
}

func NewSaveUser(repo repositories.AuthRepository, crypto providers.CryptoProvider) *SaveUser {
	return &SaveUser{repo: repo, crypto: crypto}
}

func (uc *SaveUser) Execute(ctx context.Context, input SaveUserInput) error {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))

	if name == "" {
		return models.ErrNameRequired
	}
	if email == "" {
		return models.ErrEmailRequired
	}

	existing, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return uc.create(ctx, input.ID, name, email, input.Password)
	}

	return uc.update(ctx, existing, name, email, input.Password)
}

func (uc *SaveUser) create(ctx context.Context, id, name, email, password string) error {
	if password == "" {
		return models.ErrPasswordRequired
	}
	if len(password) < 6 {
		return models.ErrPasswordTooShort
	}

	hashed, err := uc.crypto.HashPassword(ctx, password)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	user := &models.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  hashed,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return uc.repo.Create(ctx, user)
}

func (uc *SaveUser) update(ctx context.Context, existing *models.User, name, email, password string) error {
	if email != existing.Email {
		if _, err := uc.repo.FindByEmail(ctx, email); err == nil {
			return models.ErrEmailDuplicate
		}
	}

	if password != "" {
		if len(password) < 6 {
			return models.ErrPasswordTooShort
		}
		hashed, err := uc.crypto.HashPassword(ctx, password)
		if err != nil {
			return err
		}
		existing.Password = hashed
	}

	existing.Name = name
	existing.Email = email
	existing.UpdatedAt = time.Now().UTC()

	return uc.repo.Update(ctx, existing)
}
