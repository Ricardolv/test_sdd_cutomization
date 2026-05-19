package providers

import "context"

type CryptoProvider interface {
	HashPassword(ctx context.Context, password string) (string, error)
	ComparePassword(ctx context.Context, password, hash string) error
}
