package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var tokenKey string

// GenerateJWT generates a new JWT access token
func GenerateJWT(userRole string) (string, error) {
	claims := TokenClaims{
		userRole,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(tokenKey))
}

// ParseJWT parses the JWT token, checking for correct signing method.
// Returns isValid - boolean for the validation state of the token and claims data
func ParseJWT(accessToken string) (claims *TokenClaims, isValid bool, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenKey), nil
	})

	if err != nil {
		return
	} else if userClaims, ok := token.Claims.(*TokenClaims); ok {
		return userClaims, token.Valid, nil
	}
	return
}

// GetJWTKey gets the JWT secret and uses it for new token signing and old tokens verification
func GetJWTKey(key string) {
	tokenKey = key
}
