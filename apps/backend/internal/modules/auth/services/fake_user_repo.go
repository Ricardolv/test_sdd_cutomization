package services

import (
	"context"
	"sort"

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

func (f *fakeUserRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	u, exists := f.users[id]
	if !exists || !u.Active {
		return nil, models.ErrUserNotFound
	}
	return u, nil
}

func (f *fakeUserRepo) Update(ctx context.Context, u *models.User) error {
	existing, exists := f.users[u.ID]
	if !exists || !existing.Active {
		return models.ErrUserNotFound
	}
	if existing.Email != u.Email {
		if existingByEmail, duplicate := f.byEmail[u.Email]; duplicate && existingByEmail.ID != u.ID {
			return models.ErrEmailDuplicate
		}
		delete(f.byEmail, existing.Email)
	}
	f.users[u.ID] = u
	f.byEmail[u.Email] = u
	return nil
}

func (f *fakeUserRepo) Delete(ctx context.Context, id string) error {
	u, exists := f.users[id]
	if !exists || !u.Active {
		return models.ErrUserNotFound
	}
	u.Active = false
	delete(f.byEmail, u.Email)
	return nil
}

func (f *fakeUserRepo) ListPaginated(ctx context.Context, offset, limit int) ([]models.User, error) {
	var activeUsers []*models.User
	for _, u := range f.users {
		if u.Active {
			activeUsers = append(activeUsers, u)
		}
	}
	sort.Slice(activeUsers, func(i, j int) bool {
		return activeUsers[i].CreatedAt.After(activeUsers[j].CreatedAt)
	})
	if offset >= len(activeUsers) {
		return []models.User{}, nil
	}
	end := offset + limit
	if end > len(activeUsers) {
		end = len(activeUsers)
	}
	result := make([]models.User, 0, end-offset)
	for _, u := range activeUsers[offset:end] {
		result = append(result, *u)
	}
	return result, nil
}

func (f *fakeUserRepo) Count(ctx context.Context) (int, error) {
	count := 0
	for _, u := range f.users {
		if u.Active {
			count++
		}
	}
	return count, nil
}
