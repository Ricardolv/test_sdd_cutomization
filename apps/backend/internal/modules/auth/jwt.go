package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sdd-cod3r/test-app/apps/backend/internal/modules/auth/models"
)

type UserClaims struct {
	ID    string `json:"sub"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func SignUserToken(user *models.User, secret string) (string, error) {
	claims := UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(14 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
