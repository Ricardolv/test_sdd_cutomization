package models

import "errors"

var (
	ErrNameRequired     = errors.New("name is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailDuplicate   = errors.New("email must be unique")
	ErrCustomerNotFound = errors.New("customer not found")
)
