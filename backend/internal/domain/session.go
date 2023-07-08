package domain

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"
)

type Session struct {
	Username string
	Expiry   time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func SessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
