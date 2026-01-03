package authentication

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// Create token with specific email, role and expiration time
func GenerateToken(secret string, claims map[string]interface{}, duration time.Duration) (string, error) {

	jwtClaims := jwt.MapClaims{}

	for k, v := range claims {
		jwtClaims[k] = v
	}

	jwtClaims["exp"] = time.Now().Add(time.Hour * duration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(secret))
}

// Check token signature and validity return claims
func ParseTokenClaims(secret, tokenString string) (jwt.MapClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
