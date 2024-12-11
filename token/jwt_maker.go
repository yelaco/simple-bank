package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// JwtMaker is a JSON Web Token maker
type JwtMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWT maker
func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("token.NewJwtMaker: secret key must be at least %d characters", minSecretKeySize)
	}

	return &JwtMaker{secretKey}, nil
}

func (maker *JwtMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}
	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		payload,
	)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, err
	}
	return token, payload, err
}

func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("token.JwtMaker.CreateToken: %w", ErrUnexpectedSigningMethod)
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("token.JwtMaker.VerifyToken: %w", err)
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf("token.JwtMaker.VerifyToken: %w", ErrInvalidToken)
	}

	return payload, nil
}
