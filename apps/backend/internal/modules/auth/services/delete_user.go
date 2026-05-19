package services

import (
	"context"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/repositories"
)

type DeleteUserInput struct {
	ID string
}

type DeleteUser struct {
	repo repositories.AuthRepository
}

func NewDeleteUser(repo repositories.AuthRepository) *DeleteUser {
	return &DeleteUser{repo: repo}
}

func (uc *DeleteUser) Execute(ctx context.Context, input DeleteUserInput) error {
	return uc.repo.Delete(ctx, input.ID)
}
