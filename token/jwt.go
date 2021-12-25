package token

import (
	"errors"
	"time"

	"gitlab.com/canco1/canco-kit/common_error"
	"gitlab.com/canco1/canco-kit/requestinfo"

	"github.com/reddit/jwt-go"
)

type JWTAuthenticator struct {
	secretKey      string
	expirationTime time.Duration
}

func NewJWTAuthenticator(secretKey string, expirationTime time.Duration) (Authenticator, error) {
	return &JWTAuthenticator{
		secretKey,
		expirationTime,
	}, nil
}

func (a *JWTAuthenticator) Generate(info *requestinfo.Info) (*Token, error) {
	payload := NewPayload(info, a.expirationTime)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(a.secretKey))
	if err != nil {
		return nil, common_error.ErrCanNotCreateToken
	}

	return &Token{
		Token:     token,
		IssueAt:   payload.IssueAt,
		ExpiredAt: payload.ExpiredAt,
	}, nil

}

func (a *JWTAuthenticator) Verify(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, common_error.ErrInvalidToken
		}
		return []byte(a.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, common_error.ErrExpiredToken) {
			return nil, common_error.ErrExpiredToken
		}
		return nil, common_error.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, common_error.ErrInvalidToken
	}
	return payload, nil
}
