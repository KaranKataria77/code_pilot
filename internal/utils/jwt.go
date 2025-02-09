package utils

import (
	"time"

	"code_pilot/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(config.GetEnv("JWT_SECRET", "secret"))

func GenerateJWTToken(email string) (string, error) {
	expirationDate := time.Now().Add(24 * time.Hour)

	claims := &Claim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claim, error) {
	claims := &Claim{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
