package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// Maker is a generic interface for managing tokens
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

// CreateToken generates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	// Create token payload
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign token
	tkn, err := token.SignedString([]byte(maker.secretKey))
	return tkn, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunct := func(token *jwt.Token) (interface{}, error) {
		// Check token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		// Return secret key
		return []byte(maker.secretKey), nil
	}

	// Parse token
	tkn, err := jwt.ParseWithClaims(token, &Payload{}, keyFunct)
	if err != nil {
		return nil, err
	}

	// Validate token
	payload, ok := tkn.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return payload, nil
}