package token

import "time"

type Token struct {
	Token     string
	IssueAt   time.Time
	ExpiredAt time.Time
}
