package services

import (
	"context"
	"testing"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
)

func TestDeleteUser(t *testing.T) {
	repo := newFakeUserRepo()
	crypto := newFakeCryptoProvider()
	saveUser := NewSaveUser(repo, crypto)
	deleteUser := NewDeleteUser(repo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		createInput := SaveUserInput{
			ID:       "delete-me",
			Name:     "To Delete",
			Email:    "delete@example.com",
			Password: "secret123",
		}
		if err := saveUser.Execute(ctx, createInput); err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		err := deleteUser.Execute(ctx, DeleteUserInput{ID: "delete-me"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = repo.FindByID(ctx, "delete-me")
		if err != models.ErrUserNotFound {
			t.Fatalf("expected ErrUserNotFound after delete, got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		err := deleteUser.Execute(ctx, DeleteUserInput{ID: "nonexistent"})
		if err != models.ErrUserNotFound {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("double delete", func(t *testing.T) {
		createInput := SaveUserInput{
			ID:       "double-delete",
			Name:     "Double",
			Email:    "double@example.com",
			Password: "secret123",
		}
		if err := saveUser.Execute(ctx, createInput); err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		if err := deleteUser.Execute(ctx, DeleteUserInput{ID: "double-delete"}); err != nil {
			t.Fatalf("first delete failed: %v", err)
		}

		err := deleteUser.Execute(ctx, DeleteUserInput{ID: "double-delete"})
		if err != models.ErrUserNotFound {
			t.Fatalf("expected ErrUserNotFound on second delete, got %v", err)
		}
	})
}
