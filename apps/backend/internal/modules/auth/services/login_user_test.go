package services

import (
	"context"
	"testing"

	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
)

func TestLoginUser_ValidCredentials(t *testing.T) {
	repo := newFakeUserRepo()
	crypto := newFakeCryptoProvider()

	hash, _ := crypto.HashPassword(context.Background(), "secret123")
	user, _ := models.NewUser("John Doe", "john@example.com", "secret123")
	user.Password = hash
	repo.users[user.ID] = user
	repo.byEmail[user.Email] = user

	uc := NewLoginUser(repo, crypto)

	output, err := uc.Execute(context.Background(), LoginUserInput{
		Email:    "john@example.com",
		Password: "secret123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, output.ID)
	}
	if output.Name != "John Doe" {
		t.Errorf("expected name John Doe, got %s", output.Name)
	}
	if output.Email != "john@example.com" {
		t.Errorf("expected email john@example.com, got %s", output.Email)
	}
}

func TestLoginUser_EmailNotFound(t *testing.T) {
	repo := newFakeUserRepo()
	crypto := newFakeCryptoProvider()

	uc := NewLoginUser(repo, crypto)

	_, err := uc.Execute(context.Background(), LoginUserInput{
		Email:    "nonexistent@example.com",
		Password: "secret123",
	})

	if err != models.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	repo := newFakeUserRepo()
	crypto := newFakeCryptoProvider()

	hash, _ := crypto.HashPassword(context.Background(), "correct")
	user, _ := models.NewUser("John", "john@example.com", "correct")
	user.Password = hash
	repo.users[user.ID] = user
	repo.byEmail[user.Email] = user

	uc := NewLoginUser(repo, crypto)

	_, err := uc.Execute(context.Background(), LoginUserInput{
		Email:    "john@example.com",
		Password: "wrong",
	})

	if err != models.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginUser_EmptyEmail(t *testing.T) {
	repo := newFakeUserRepo()
	crypto := newFakeCryptoProvider()
	uc := NewLoginUser(repo, crypto)

	_, err := uc.Execute(context.Background(), LoginUserInput{
		Email:    "",
		Password: "secret123",
	})

	if err == nil {
		t.Fatal("expected error for empty email")
	}
	if err.Error() != "email is required" {
		t.Fatalf("expected 'email is required', got %v", err)
	}
}

func TestLoginUser_InvalidEmail(t *testing.T) {
	repo := newFakeUserRepo()
	crypto := newFakeCryptoProvider()
	uc := NewLoginUser(repo, crypto)

	_, err := uc.Execute(context.Background(), LoginUserInput{
		Email:    "not-an-email",
		Password: "secret123",
	})

	if err == nil {
		t.Fatal("expected error for invalid email")
	}
	if err.Error() != "email is invalid" {
		t.Fatalf("expected 'email is invalid', got %v", err)
	}
}

func TestLoginUser_EmptyPassword(t *testing.T) {
	repo := newFakeUserRepo()
	crypto := newFakeCryptoProvider()
	uc := NewLoginUser(repo, crypto)

	_, err := uc.Execute(context.Background(), LoginUserInput{
		Email:    "john@example.com",
		Password: "",
	})

	if err == nil {
		t.Fatal("expected error for empty password")
	}
	if err.Error() != "password is required" {
		t.Fatalf("expected 'password is required', got %v", err)
	}
}
