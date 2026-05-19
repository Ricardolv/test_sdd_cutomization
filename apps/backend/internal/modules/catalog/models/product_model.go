package models

import (
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	StatusActive   ProductStatus = "active"
	StatusInactive ProductStatus = "inactive"
	StatusDraft    ProductStatus = "draft"
)

func ValidProductStatuses() []ProductStatus {
	return []ProductStatus{StatusActive, StatusInactive, StatusDraft}
}

func (s ProductStatus) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusDraft:
		return true
	}
	return false
}

type Product struct {
	ID              string
	Name            string
	Description     string
	Price           float64
	Status          ProductStatus
	AvailableOnline bool
	Featured        bool
	AllowsPreOrder  bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewProduct(name, description string, price float64, status ProductStatus, availableOnline, featured, allowsPreOrder bool) (*Product, error) {
	p := &Product{
		ID:              uuid.New().String(),
		Name:            strings.TrimSpace(name),
		Description:     strings.TrimSpace(description),
		Price:           price,
		Status:          status,
		AvailableOnline: availableOnline,
		Featured:        featured,
		AllowsPreOrder:  allowsPreOrder,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrNameRequired
	}
	if len(p.Name) < 2 {
		return ErrNameTooShort
	}
	if len(p.Name) > 120 {
		return ErrNameTooLong
	}
	if len(p.Description) > 500 {
		return ErrDescriptionTooLong
	}
	if p.Price < 0 {
		return ErrPriceNegative
	}
	rounded := math.Round(p.Price*100) / 100
	if math.Abs(p.Price-rounded) > 1e-9 {
		return ErrPricePrecision
	}
	if !p.Status.IsValid() {
		return ErrStatusInvalid
	}
	return nil
}

func (p *Product) Update(name, description string, price float64, status ProductStatus, availableOnline, featured, allowsPreOrder bool) error {
	p.Name = strings.TrimSpace(name)
	p.Description = strings.TrimSpace(description)
	p.Price = price
	p.Status = status
	p.AvailableOnline = availableOnline
	p.Featured = featured
	p.AllowsPreOrder = allowsPreOrder
	p.UpdatedAt = time.Now().UTC()

	return p.Validate()
}
