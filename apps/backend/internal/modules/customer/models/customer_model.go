package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        string
	Name      string
	Email     string
	Phone     string
	Notes     string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCustomer(name, email, phone, notes string) (*Customer, error) {
	c := &Customer{
		ID:        uuid.New().String(),
		Name:      strings.TrimSpace(name),
		Email:     strings.TrimSpace(email),
		Phone:     strings.TrimSpace(phone),
		Notes:     strings.TrimSpace(notes),
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Customer) Validate() error {
	if c.Name == "" {
		return ErrNameRequired
	}
	if c.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (c *Customer) Update(name, email, phone, notes string) error {
	c.Name = strings.TrimSpace(name)
	c.Email = strings.TrimSpace(email)
	c.Phone = strings.TrimSpace(phone)
	c.Notes = strings.TrimSpace(notes)
	c.UpdatedAt = time.Now().UTC()

	return c.Validate()
}

func (c *Customer) Deactivate() {
	c.Active = false
	c.UpdatedAt = time.Now().UTC()
}

func IsEmailDuplicateError(err error) bool {
	return errors.Is(err, ErrEmailDuplicate)
}
