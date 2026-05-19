package services

import (
	"context"
	"testing"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
)

func TestSaveUser_Create(t *testing.T) {
	repo := newFakeUserRepo()
	crypto := newFakeCryptoProvider()
	uc := NewSaveUser(repo, crypto)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		input := SaveUserInput{
			ID:       "test-id-1",
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "secret123",
		}
		err := uc.Execute(ctx, input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		user, err := repo.FindByID(ctx, "test-id-1")
		if err != nil {
			t.Fatalf("expected user to exist, got %v", err)
		}
		if user.Name != "Test User" {
			t.Errorf("expected name Test User, got %s", user.Name)
		}
		if user.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", user.Email)
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		input := SaveUserInput{
			ID:       "test-id-2",
			Name:     "Another User",
			Email:    "test@example.com",
			Password: "secret123",
		}
		err := uc.Execute(ctx, input)
		if err != models.ErrEmailDuplicate {
			t.Fatalf("expected ErrEmailDuplicate, got %v", err)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		input := SaveUserInput{
			ID:       "test-id-3",
			Name:     "",
			Email:    "new@example.com",
			Password: "secret123",
		}
		err := uc.Execute(ctx, input)
		if err != models.ErrNameRequired {
			t.Fatalf("expected ErrNameRequired, got %v", err)
		}
	})

	t.Run("empty email", func(t *testing.T) {
		input := SaveUserInput{
			ID:       "test-id-4",
			Name:     "Test User",
			Email:    "",
			Password: "secret123",
		}
		err := uc.Execute(ctx, input)
		if err != models.ErrEmailRequired {
			t.Fatalf("expected ErrEmailRequired, got %v", err)
		}
	})

	t.Run("empty password on create", func(t *testing.T) {
		input := SaveUserInput{
			ID:       "test-id-5",
			Name:     "Test User",
			Email:    "nopass@example.com",
			Password: "",
		}
		err := uc.Execute(ctx, input)
		if err != models.ErrPasswordRequired {
			t.Fatalf("expected ErrPasswordRequired, got %v", err)
		}
	})

	t.Run("password too short", func(t *testing.T) {
		input := SaveUserInput{
			ID:       "test-id-6",
			Name:     "Test User",
			Email:    "short@example.com",
			Password: "12345",
		}
		err := uc.Execute(ctx, input)
		if err != models.ErrPasswordTooShort {
			t.Fatalf("expected ErrPasswordTooShort, got %v", err)
		}
	})
}

func TestSaveUser_Update(t *testing.T) {
	repo := newFakeUserRepo()
	crypto := newFakeCryptoProvider()
	uc := NewSaveUser(repo, crypto)
	ctx := context.Background()

	existingID := "existing-id"
	input := SaveUserInput{
		ID:       existingID,
		Name:     "Original",
		Email:    "original@example.com",
		Password: "secret123",
	}
	if err := uc.Execute(ctx, input); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Run("success with new password", func(t *testing.T) {
		updateInput := SaveUserInput{
			ID:       existingID,
			Name:     "Updated Name",
			Email:    "updated@example.com",
			Password: "newpass123",
		}
		err := uc.Execute(ctx, updateInput)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		user, err := repo.FindByID(ctx, existingID)
		if err != nil {
			t.Fatalf("expected user to exist, got %v", err)
		}
		if user.Name != "Updated Name" {
			t.Errorf("expected name Updated Name, got %s", user.Name)
		}
		if user.Email != "updated@example.com" {
			t.Errorf("expected email updated@example.com, got %s", user.Email)
		}
	})

	t.Run("success keeping password", func(t *testing.T) {
		updateInput := SaveUserInput{
			ID:       existingID,
			Name:     "No Password Change",
			Email:    "nopasschange@example.com",
			Password: "",
		}
		err := uc.Execute(ctx, updateInput)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		user, err := repo.FindByID(ctx, existingID)
		if err != nil {
			t.Fatalf("expected user to exist, got %v", err)
		}
		if user.Name != "No Password Change" {
			t.Errorf("expected name No Password Change, got %s", user.Name)
		}
	})

	t.Run("password too short on update", func(t *testing.T) {
		updateInput := SaveUserInput{
			ID:       existingID,
			Name:     "Bad Password",
			Email:    "badpass@example.com",
			Password: "123",
		}
		err := uc.Execute(ctx, updateInput)
		if err != models.ErrPasswordTooShort {
			t.Fatalf("expected ErrPasswordTooShort, got %v", err)
		}
	})

	t.Run("duplicate email on update", func(t *testing.T) {
		anotherID := "another-id"
		anotherInput := SaveUserInput{
			ID:       anotherID,
			Name:     "Another",
			Email:    "another@example.com",
			Password: "secret123",
		}
		if err := uc.Execute(ctx, anotherInput); err != nil {
			t.Fatalf("failed to create another user: %v", err)
		}

		updateInput := SaveUserInput{
			ID:       existingID,
			Name:     "Steal Email",
			Email:    "another@example.com",
			Password: "",
		}
		err := uc.Execute(ctx, updateInput)
		if err != models.ErrEmailDuplicate {
			t.Fatalf("expected ErrEmailDuplicate, got %v", err)
		}
	})
}
