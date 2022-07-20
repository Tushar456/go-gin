package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type JWTToken struct {
	secret string
}

func NewJWTToken(secretkey string) (Token, error) {

	if len(secretkey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTToken{secretkey}, nil
}

func (token *JWTToken) CreateToken(username string, duration time.Duration) (string, *Payload, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := jwtToken.SignedString([]byte(token.secret))
	if err != nil {
		return "", payload, err
	}
	return tokenString, payload, nil
}

func (token *JWTToken) VerifyToken(tokenString string) (*Payload, error) {

	keyFunc := func(jwtToken *jwt.Token) (interface{}, error) {
		_, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(token.secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil

}
