package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenClaims struct {
	Role     string `json:"role"`
	Username string `json:"username"`
	Id       int    `json:"id"`
	jwt.RegisteredClaims
}

type GenerateJWTParams struct {
	Role       string
	Username   string
	Id         int
	Expiration time.Duration
}

var tokenKey string

// GenerateJWT generates a new JWT access token
func GenerateJWT(params GenerateJWTParams) (string, error) {
	var expirationTime time.Duration

	if params.Expiration != 0 {
		expirationTime = params.Expiration
	} else {
		expirationTime = 720 * time.Hour // 30 days
	}

	claims := TokenClaims{
		params.Role,
		params.Username,
		params.Id,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
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
