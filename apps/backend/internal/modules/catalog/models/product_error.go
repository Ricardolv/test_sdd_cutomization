package models

import "errors"

var (
	ErrNameRequired          = errors.New("name is required")
	ErrNameTooShort          = errors.New("name must be at least 2 characters")
	ErrNameTooLong           = errors.New("name must be at most 120 characters")
	ErrDescriptionTooLong    = errors.New("description must be at most 500 characters")
	ErrPriceRequired         = errors.New("price is required")
	ErrPriceNegative         = errors.New("price must be non-negative")
	ErrPricePrecision        = errors.New("price must have at most 2 decimal places")
	ErrStatusRequired        = errors.New("status is required")
	ErrStatusInvalid         = errors.New("status must be active, inactive, or draft")
	ErrProductNotFound       = errors.New("product not found")
)
