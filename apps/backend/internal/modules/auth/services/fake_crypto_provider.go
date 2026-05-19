package services

import (
	"context"
	"errors"
)

type fakeCryptoProvider struct {
	passwords map[string]string // plain -> hash
}

func newFakeCryptoProvider() *fakeCryptoProvider {
	return &fakeCryptoProvider{
		passwords: make(map[string]string),
	}
}

func (f *fakeCryptoProvider) HashPassword(_ context.Context, password string) (string, error) {
	hash := "hashed_" + password
	f.passwords[password] = hash
	return hash, nil
}

func (f *fakeCryptoProvider) ComparePassword(_ context.Context, password, hash string) error {
	expected, ok := f.passwords[password]
	if !ok {
		return errors.New("password mismatch")
	}
	if expected != hash {
		return errors.New("password mismatch")
	}
	return nil
}
