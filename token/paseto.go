package token

import (
	"time"

	"gitlab.com/canco1/canco-kit/common_error"
	"gitlab.com/canco1/canco-kit/requestinfo"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoAuthenticator struct {
	paseto         *paseto.V2
	symmetricKey   []byte
	expirationTime time.Duration
}

func NewPasetoAuthenticator(symmetricKey string, expirationTime time.Duration) (Authenticator, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, common_error.ErrInvalidKeySize
	}
	return &PasetoAuthenticator{
		paseto:         paseto.NewV2(),
		symmetricKey:   []byte(symmetricKey),
		expirationTime: expirationTime,
	}, nil
}

func (a *PasetoAuthenticator) Generate(info *requestinfo.Info) (*Token, error) {
	payload := NewPayload(info, a.expirationTime)
	token, err := a.paseto.Encrypt(a.symmetricKey, payload, nil)

	if err != nil {
		return nil, common_error.ErrCanNotCreateToken
	}

	return &Token{
		Token:     token,
		ExpiredAt: payload.ExpiredAt,
		IssueAt:   payload.IssueAt,
	}, nil
}

func (a *PasetoAuthenticator) Verify(token string) (*Payload, error) {
	payload := &Payload{}

	if err := a.paseto.Decrypt(token, a.symmetricKey, payload, nil); err != nil {
		return nil, common_error.ErrInvalidToken
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
