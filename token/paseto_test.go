package token

import (
	"testing"
	"time"

	"gitlab.com/canco1/canco-kit/common_error"
	"gitlab.com/canco1/canco-kit/utils"

	"github.com/stretchr/testify/assert"
)

func Test_Paseto(t *testing.T) {
	authenticator, err := NewPasetoAuthenticator(utils.RandStringBytes(32), time.Minute)

	assert.NoError(t, err)

	token, err := authenticator.Generate(info)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.NotEmpty(t, token.ExpiredAt)

	payload, err := authenticator.Verify(token.Token)

	assert.NoError(t, err)
	assert.NotNil(t, payload)
	assert.Equal(t, payload.UserID, info.UserID)
	assert.Equal(t, payload.UserName, info.UserName)
	assert.Equal(t, payload.Ip, info.Ip)
	assert.Equal(t, payload.UserAgent, info.UserAgent)
	assert.WithinDuration(t, payload.ExpiredAt, token.ExpiredAt, time.Second)
	assert.WithinDuration(t, payload.IssueAt, token.IssueAt, time.Second)
}

func Test_Paseto_ExpiredToken(t *testing.T) {
	authenticator, err := NewPasetoAuthenticator(utils.RandStringBytes(32), -time.Minute)

	assert.NoError(t, err)

	token, err := authenticator.Generate(info)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.NotEmpty(t, token.ExpiredAt)

	payload, err := authenticator.Verify(token.Token)
	assert.Error(t, err)
	assert.Equal(t, err, common_error.ErrExpiredToken)
	assert.Nil(t, payload)
}
