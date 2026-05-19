package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, email, password string) (*User, error) {
	u := &User{
		ID:        uuid.New().String(),
		Name:      strings.TrimSpace(name),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Password:  password,
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Validate() error {
	var errs []error

	if u.Name == "" {
		errs = append(errs, ErrNameRequired)
	}
	if u.Email == "" {
		errs = append(errs, ErrEmailRequired)
	}
	if u.Password == "" {
		errs = append(errs, ErrPasswordRequired)
	} else if len(u.Password) < 6 {
		errs = append(errs, ErrPasswordTooShort)
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}

	return nil
}

type ValidationError struct {
	Errors []error
}

func (e *ValidationError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return "multiple validation errors"
}

func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

func GetValidationErrors(err error) []error {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve.Errors
	}
	return nil
}

var (
	ErrNameRequired     = errors.New("name is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailDuplicate   = errors.New("email must be unique")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
	ErrUserNotFound     = errors.New("user not found")
)
