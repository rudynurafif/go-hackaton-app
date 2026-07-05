package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims — isi token. Menggantikan tabel session milik Better Auth:
// identitas dibawa di dalam token (stateless), tidak disimpan di database.
type JWTClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT bertanda tangan HS256 yang memuat id user
// (subject) dan role, berlaku selama `expires`.
func GenerateToken(userID, role, secret string, expires time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expires)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken memverifikasi tanda tangan + masa berlaku token dan
// mengembalikan claims-nya.
func ParseToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(t *jwt.Token) (any, error) { return []byte(secret), nil },
		// Tolak token dengan algoritma lain (mis. "none") — hanya HS256.
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
