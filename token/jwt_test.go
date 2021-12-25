package token

import (
	"testing"
	"time"

	"gitlab.com/canco1/canco-kit/common_error"
	"gitlab.com/canco1/canco-kit/requestinfo"
	"gitlab.com/canco1/canco-kit/utils"

	"github.com/bxcodec/faker/v3"
	"github.com/reddit/jwt-go"
	"github.com/stretchr/testify/assert"
)

var info = &requestinfo.Info{
	UserID:    faker.UUIDDigit(),
	UserName:  faker.Username(),
	Ip:        faker.IPv4(),
	UserAgent: faker.MacAddress(),
}

func Test_JWT(t *testing.T) {
	authenticator, err := NewJWTAuthenticator(utils.RandStringBytes(32), time.Minute)

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

func Test_JWT_ExpiredToken(t *testing.T) {
	authenticator, err := NewJWTAuthenticator(utils.RandStringBytes(32), -time.Minute)

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

func Test_JWT_InvalidToken(t *testing.T) {
	authenticator, err := NewJWTAuthenticator(utils.RandStringBytes(32), time.Minute)

	assert.NoError(t, err)

	token, err := authenticator.Generate(info)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	assert.NotEmpty(t, token.ExpiredAt)

	payload := NewPayload(info, time.Minute)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	tknString, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	p, err := authenticator.Verify(tknString)
	assert.Error(t, err)
	assert.Equal(t, err, common_error.ErrInvalidToken)
	assert.Nil(t, p)
}
