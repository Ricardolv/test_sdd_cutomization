package providers

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type bcryptCryptoProvider struct{}

func NewBcryptCryptoProvider() CryptoProvider {
	return &bcryptCryptoProvider{}
}

func (p *bcryptCryptoProvider) HashPassword(_ context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (p *bcryptCryptoProvider) ComparePassword(_ context.Context, password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
