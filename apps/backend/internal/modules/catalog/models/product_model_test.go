package models

import (
	"strings"
	"testing"
)

func TestNewProduct_Success(t *testing.T) {
	p, err := NewProduct("Widget", "A useful widget", 19.99, StatusActive, true, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.ID == "" {
		t.Error("expected ID to be generated")
	}
	if p.Name != "Widget" {
		t.Errorf("expected name 'Widget', got %s", p.Name)
	}
	if p.Description != "A useful widget" {
		t.Errorf("expected description 'A useful widget', got %s", p.Description)
	}
	if p.Price != 19.99 {
		t.Errorf("expected price 19.99, got %f", p.Price)
	}
	if p.Status != StatusActive {
		t.Errorf("expected status active, got %s", p.Status)
	}
	if !p.AvailableOnline {
		t.Error("expected AvailableOnline to be true")
	}
	if p.Featured {
		t.Error("expected Featured to be false")
	}
	if p.AllowsPreOrder {
		t.Error("expected AllowsPreOrder to be false")
	}
	if p.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if p.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestNewProduct_TrimsWhitespace(t *testing.T) {
	p, err := NewProduct("  Widget  ", "  desc  ", 10.00, StatusDraft, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Name != "Widget" {
		t.Errorf("expected trimmed name, got %q", p.Name)
	}
	if p.Description != "desc" {
		t.Errorf("expected trimmed description, got %q", p.Description)
	}
}

func TestNewProduct_EmptyDescription(t *testing.T) {
	p, err := NewProduct("Widget", "", 10.00, StatusActive, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Description != "" {
		t.Errorf("expected empty description, got %q", p.Description)
	}
}

func TestNewProduct_MissingName(t *testing.T) {
	_, err := NewProduct("", "desc", 10.00, StatusActive, false, false, false)
	if err != ErrNameRequired {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestNewProduct_NameTooShort(t *testing.T) {
	_, err := NewProduct("A", "desc", 10.00, StatusActive, false, false, false)
	if err != ErrNameTooShort {
		t.Fatalf("expected ErrNameTooShort, got %v", err)
	}
}

func TestNewProduct_NameTooLong(t *testing.T) {
	longName := strings.Repeat("A", 121)
	_, err := NewProduct(longName, "desc", 10.00, StatusActive, false, false, false)
	if err != ErrNameTooLong {
		t.Fatalf("expected ErrNameTooLong, got %v", err)
	}
}

func TestNewProduct_DescriptionTooLong(t *testing.T) {
	longDesc := strings.Repeat("A", 501)
	_, err := NewProduct("Widget", longDesc, 10.00, StatusActive, false, false, false)
	if err != ErrDescriptionTooLong {
		t.Fatalf("expected ErrDescriptionTooLong, got %v", err)
	}
}

func TestNewProduct_NegativePrice(t *testing.T) {
	_, err := NewProduct("Widget", "desc", -1.00, StatusActive, false, false, false)
	if err != ErrPriceNegative {
		t.Fatalf("expected ErrPriceNegative, got %v", err)
	}
}

func TestNewProduct_ZeroPrice(t *testing.T) {
	p, err := NewProduct("Widget", "desc", 0, StatusActive, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Price != 0 {
		t.Errorf("expected price 0, got %f", p.Price)
	}
}

func TestNewProduct_PriceTooManyDecimals(t *testing.T) {
	_, err := NewProduct("Widget", "desc", 10.999, StatusActive, false, false, false)
	if err != ErrPricePrecision {
		t.Fatalf("expected ErrPricePrecision, got %v", err)
	}
}

func TestNewProduct_PriceTwoDecimals(t *testing.T) {
	p, err := NewProduct("Widget", "desc", 10.99, StatusActive, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Price != 10.99 {
		t.Errorf("expected price 10.99, got %f", p.Price)
	}
}

func TestNewProduct_InvalidStatus(t *testing.T) {
	_, err := NewProduct("Widget", "desc", 10.00, ProductStatus("unknown"), false, false, false)
	if err != ErrStatusInvalid {
		t.Fatalf("expected ErrStatusInvalid, got %v", err)
	}
}

func TestNewProduct_AllStatuses(t *testing.T) {
	statuses := []ProductStatus{StatusActive, StatusInactive, StatusDraft}
	for _, s := range statuses {
		p, err := NewProduct("Widget", "desc", 10.00, s, false, false, false)
		if err != nil {
			t.Fatalf("unexpected error for status %s: %v", s, err)
		}
		if p.Status != s {
			t.Errorf("expected status %s, got %s", s, p.Status)
		}
	}
}

func TestProductUpdate_Success(t *testing.T) {
	p, _ := NewProduct("Widget", "desc", 10.00, StatusActive, false, false, false)

	err := p.Update("Gadget", "updated desc", 25.50, StatusInactive, true, true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Name != "Gadget" {
		t.Errorf("expected name 'Gadget', got %s", p.Name)
	}
	if p.Description != "updated desc" {
		t.Errorf("expected description 'updated desc', got %s", p.Description)
	}
	if p.Price != 25.50 {
		t.Errorf("expected price 25.50, got %f", p.Price)
	}
	if p.Status != StatusInactive {
		t.Errorf("expected status inactive, got %s", p.Status)
	}
	if !p.AvailableOnline {
		t.Error("expected AvailableOnline to be true")
	}
	if !p.Featured {
		t.Error("expected Featured to be true")
	}
	if !p.AllowsPreOrder {
		t.Error("expected AllowsPreOrder to be true")
	}
}

func TestProductUpdate_ValidationFailure(t *testing.T) {
	p, _ := NewProduct("Widget", "desc", 10.00, StatusActive, false, false, false)

	err := p.Update("", "desc", 10.00, StatusActive, false, false, false)
	if err != ErrNameRequired {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestProductUpdate_TrimsWhitespace(t *testing.T) {
	p, _ := NewProduct("Widget", "desc", 10.00, StatusActive, false, false, false)

	err := p.Update("  Gadget  ", "  updated  ", 15.00, StatusDraft, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Name != "Gadget" {
		t.Errorf("expected trimmed name, got %q", p.Name)
	}
	if p.Description != "updated" {
		t.Errorf("expected trimmed description, got %q", p.Description)
	}
}

func TestProductUpdate_PreservesID(t *testing.T) {
	p, _ := NewProduct("Widget", "desc", 10.00, StatusActive, false, false, false)
	originalID := p.ID

	_ = p.Update("Gadget", "updated", 15.00, StatusDraft, false, false, false)

	if p.ID != originalID {
		t.Errorf("expected ID to remain %s, got %s", originalID, p.ID)
	}
}

func TestProductStatus_IsValid(t *testing.T) {
	if !StatusActive.IsValid() {
		t.Error("expected active to be valid")
	}
	if !StatusInactive.IsValid() {
		t.Error("expected inactive to be valid")
	}
	if !StatusDraft.IsValid() {
		t.Error("expected draft to be valid")
	}
	if ProductStatus("unknown").IsValid() {
		t.Error("expected unknown to be invalid")
	}
	if ProductStatus("").IsValid() {
		t.Error("expected empty to be invalid")
	}
}

func TestValidProductStatuses(t *testing.T) {
	statuses := ValidProductStatuses()
	if len(statuses) != 3 {
		t.Fatalf("expected 3 statuses, got %d", len(statuses))
	}
}

func TestNewProduct_UniqueIDs(t *testing.T) {
	p1, _ := NewProduct("Widget", "desc", 10.00, StatusActive, false, false, false)
	p2, _ := NewProduct("Gadget", "desc", 20.00, StatusDraft, false, false, false)

	if p1.ID == p2.ID {
		t.Error("expected different IDs for different products")
	}
}

func TestNewProduct_NameExactly2Chars(t *testing.T) {
	p, err := NewProduct("AB", "desc", 10.00, StatusActive, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Name != "AB" {
		t.Errorf("expected name 'AB', got %s", p.Name)
	}
}

func TestNewProduct_NameExactly120Chars(t *testing.T) {
	name := strings.Repeat("A", 120)
	p, err := NewProduct(name, "desc", 10.00, StatusActive, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(p.Name) != 120 {
		t.Errorf("expected name length 120, got %d", len(p.Name))
	}
}

func TestNewProduct_DescriptionExactly500Chars(t *testing.T) {
	desc := strings.Repeat("A", 500)
	p, err := NewProduct("Widget", desc, 10.00, StatusActive, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(p.Description) != 500 {
		t.Errorf("expected description length 500, got %d", len(p.Description))
	}
}

func TestProductValidate_MutatesBeforeValidate(t *testing.T) {
	p, _ := NewProduct("Widget", "desc", 10.00, StatusActive, false, false, false)

	_ = p.Update("", "desc", 10.00, StatusActive, false, false, false)

	if p.Name != "" {
		t.Errorf("expected name to be empty after update (mutates before validate), got %s", p.Name)
	}
}
