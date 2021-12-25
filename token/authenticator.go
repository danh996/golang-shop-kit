package token

import "gitlab.com/canco1/canco-kit/requestinfo"

type Authenticator interface {
	Generate(payload *requestinfo.Info) (*Token, error)
	Verify(token string) (*Payload, error)
}
