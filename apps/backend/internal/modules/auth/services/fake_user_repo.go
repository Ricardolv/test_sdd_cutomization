package services

import (
	"context"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
)

type fakeUserRepo struct {
	users   map[string]*models.User
	byEmail map[string]*models.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		users:   make(map[string]*models.User),
		byEmail: make(map[string]*models.User),
	}
}

func (f *fakeUserRepo) Create(ctx context.Context, u *models.User) error {
	if _, exists := f.byEmail[u.Email]; exists {
		return models.ErrEmailDuplicate
	}
	f.users[u.ID] = u
	f.byEmail[u.Email] = u
	return nil
}

func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u, exists := f.byEmail[email]
	if !exists {
		return nil, models.ErrUserNotFound
	}
	return u, nil
}
