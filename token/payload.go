package token

import (
	"encoding/json"
	"time"

	"gitlab.com/canco1/canco-kit/common_error"
	"gitlab.com/canco1/canco-kit/requestinfo"
)

type Payload struct {
	requestinfo.Info
	ExpiredAt time.Time `json:"expired_at"`
	IssueAt   time.Time `json:"issue_at"`
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return common_error.ErrExpiredToken
	}
	return nil
}

func (p *Payload) ParseToJSONString() (string, error) {
	m, err := json.Marshal(p)
	if err != nil {
		return "", nil
	}
	return string(m), nil
}

func NewPayload(info *requestinfo.Info, expirationTime time.Duration) *Payload {
	t := time.Now()
	expiredAt := t.Add(expirationTime)
	return &Payload{
		Info:      *info,
		IssueAt:   t,
		ExpiredAt: expiredAt,
	}
}

func ParseJSONStringToPayload(str string) (*Payload, error) {
	payload := &Payload{}
	if err := json.Unmarshal([]byte(str), payload); err != nil {
		return nil, err
	}
	return payload, nil
}
