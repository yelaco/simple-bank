package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	subject  = "simple-bank"
	audience = "user"
)

// Payload contains the payload data of the token
type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return &Payload{}, err
	}

	return &Payload{
		Id:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: p.ExpiredAt,
	}, nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: p.IssuedAt,
	}, nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (p *Payload) GetIssuer() (string, error) {
	return p.Username, nil
}

func (p *Payload) GetSubject() (string, error) {
	return subject, nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{audience}, nil
}
