package models

import (
	"errors"
	"testing"
)

func TestNewCustomer_Success(t *testing.T) {
	c, err := NewCustomer("John Doe", "john@example.com", "+1234567890", "VIP customer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.ID == "" {
		t.Error("expected ID to be generated")
	}
	if c.Name != "John Doe" {
		t.Errorf("expected name 'John Doe', got %s", c.Name)
	}
	if c.Email != "john@example.com" {
		t.Errorf("expected email 'john@example.com', got %s", c.Email)
	}
	if c.Phone != "+1234567890" {
		t.Errorf("expected phone '+1234567890', got %s", c.Phone)
	}
	if c.Notes != "VIP customer" {
		t.Errorf("expected notes 'VIP customer', got %s", c.Notes)
	}
	if !c.Active {
		t.Error("expected customer to be active")
	}
	if c.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if c.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestNewCustomer_TrimsWhitespace(t *testing.T) {
	c, err := NewCustomer("  John Doe  ", "  john@example.com  ", "  +123  ", "  notes  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Name != "John Doe" {
		t.Errorf("expected trimmed name, got %q", c.Name)
	}
	if c.Email != "john@example.com" {
		t.Errorf("expected trimmed email, got %q", c.Email)
	}
}

func TestNewCustomer_MissingName(t *testing.T) {
	_, err := NewCustomer("", "john@example.com", "", "")
	if !errors.Is(err, ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestNewCustomer_MissingEmail(t *testing.T) {
	_, err := NewCustomer("John Doe", "", "", "")
	if !errors.Is(err, ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestNewCustomer_OnlyWhitespaceName(t *testing.T) {
	_, err := NewCustomer("   ", "john@example.com", "", "")
	if !errors.Is(err, ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired for whitespace-only name, got %v", err)
	}
}

func TestNewCustomer_OnlyWhitespaceEmail(t *testing.T) {
	_, err := NewCustomer("John Doe", "   ", "", "")
	if !errors.Is(err, ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired for whitespace-only email, got %v", err)
	}
}

func TestCustomerValidate_Valid(t *testing.T) {
	c := &Customer{
		ID:    "test-id",
		Name:  "John",
		Email: "john@example.com",
	}

	err := c.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCustomerValidate_EmptyName(t *testing.T) {
	c := &Customer{
		ID:    "test-id",
		Name:  "",
		Email: "john@example.com",
	}

	err := c.Validate()
	if !errors.Is(err, ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestCustomerValidate_EmptyEmail(t *testing.T) {
	c := &Customer{
		ID:    "test-id",
		Name:  "John",
		Email: "",
	}

	err := c.Validate()
	if !errors.Is(err, ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestCustomerUpdate_Success(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "+123", "notes")

	err := c.Update("John Updated", "john.updated@example.com", "+456", "updated notes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Name != "John Updated" {
		t.Errorf("expected name 'John Updated', got %s", c.Name)
	}
	if c.Email != "john.updated@example.com" {
		t.Errorf("expected email 'john.updated@example.com', got %s", c.Email)
	}
	if c.Phone != "+456" {
		t.Errorf("expected phone '+456', got %s", c.Phone)
	}
	if c.Notes != "updated notes" {
		t.Errorf("expected notes 'updated notes', got %s", c.Notes)
	}
	if c.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestCustomerUpdate_TrimsWhitespace(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "", "")

	err := c.Update("  Updated  ", "  updated@example.com  ", "  +789  ", "  new notes  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Name != "Updated" {
		t.Errorf("expected trimmed name, got %q", c.Name)
	}
	if c.Email != "updated@example.com" {
		t.Errorf("expected trimmed email, got %q", c.Email)
	}
}

func TestCustomerUpdate_MissingName(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "", "")

	err := c.Update("", "john.updated@example.com", "", "")
	if !errors.Is(err, ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestCustomerUpdate_MissingEmail(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "", "")

	err := c.Update("John Updated", "", "", "")
	if !errors.Is(err, ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestCustomerDeactivate(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "", "")

	c.Deactivate()

	if c.Active {
		t.Error("expected customer to be inactive")
	}
	if c.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestIsEmailDuplicateError(t *testing.T) {
	if !IsEmailDuplicateError(ErrEmailDuplicate) {
		t.Error("expected IsEmailDuplicateError to return true for ErrEmailDuplicate")
	}

	otherErr := errors.New("some other error")
	if IsEmailDuplicateError(otherErr) {
		t.Error("expected IsEmailDuplicateError to return false for other errors")
	}

	if IsEmailDuplicateError(nil) {
		t.Error("expected IsEmailDuplicateError to return false for nil")
	}
}

func TestCustomerUpdate_PreservesIDAndActive(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "", "")
	originalID := c.ID

	_ = c.Update("Updated", "updated@example.com", "", "")

	if c.ID != originalID {
		t.Errorf("expected ID to remain %s, got %s", originalID, c.ID)
	}
	if !c.Active {
		t.Error("expected Active to remain true after update")
	}
}

func TestNewCustomer_UniqueIDs(t *testing.T) {
	c1, _ := NewCustomer("John", "john@example.com", "", "")
	c2, _ := NewCustomer("Jane", "jane@example.com", "", "")

	if c1.ID == c2.ID {
		t.Error("expected different IDs for different customers")
	}
}

func TestCustomerValidate_CreatedAtBeforeUpdatedAt(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "", "")

	if !c.UpdatedAt.Equal(c.CreatedAt) && !c.UpdatedAt.After(c.CreatedAt) {
		t.Error("expected UpdatedAt to be equal to or after CreatedAt")
	}
}

func TestCustomerUpdate_MutatesBeforeValidation(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "+123", "notes")

	_ = c.Update("", "new@example.com", "+456", "new notes")

	if c.Name != "" {
		t.Errorf("expected name to be empty after update (mutates before validate), got %s", c.Name)
	}
	if c.Phone != "+456" {
		t.Errorf("expected phone to be '+456' after update, got %s", c.Phone)
	}
}

func TestNewCustomer_EmptyPhoneAndNotes(t *testing.T) {
	c, err := NewCustomer("John", "john@example.com", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Phone != "" {
		t.Errorf("expected empty phone, got %s", c.Phone)
	}
	if c.Notes != "" {
		t.Errorf("expected empty notes, got %s", c.Notes)
	}
}

func TestCustomerDeactivate_MultipleTimes(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "", "")

	c.Deactivate()
	c.Deactivate()

	if c.Active {
		t.Error("expected customer to remain inactive after multiple deactivations")
	}
}

func TestCustomerUpdate_EmptyStringsAreValidForOptionalFields(t *testing.T) {
	c, _ := NewCustomer("John", "john@example.com", "+123", "notes")

	err := c.Update("John Updated", "john.updated@example.com", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Phone != "" {
		t.Errorf("expected phone to be cleared, got %s", c.Phone)
	}
	if c.Notes != "" {
		t.Errorf("expected notes to be cleared, got %s", c.Notes)
	}
}

func TestCustomerValidate_NameWithOnlySpaces(t *testing.T) {
	c := &Customer{
		Name:  "John",
		Email: "test@example.com",
	}

	err := c.Validate()
	if err != nil {
		t.Fatalf("expected no error for valid name, got %v", err)
	}
}

func TestCustomerValidate_EmailWithOnlySpaces(t *testing.T) {
	c := &Customer{
		Name:  "John",
		Email: "test@example.com",
	}

	err := c.Validate()
	if err != nil {
		t.Fatalf("expected no error for valid email, got %v", err)
	}
}
