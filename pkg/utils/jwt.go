package utils

import (
	"time"

	"github.com/Beeram12/college-appointment-system/config"
	"github.com/golang-jwt/jwt/v5"
)

// Load config
var cfg = config.JwtLoadConfig()
var jwtSecret = []byte(cfg.JwtSecret)

type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Generate tokens
func GenerateToken(username, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := CustomClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "college-appointment-system",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Validating Tokens
func ValidatingTokens(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
