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

var sessions = make(map[string]Session)

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func GetSessions() map[string]Session {
	return sessions
}

func SetSession(session Session) string {
	sessionId := sessionId()

	sessions[sessionId] = session

	return sessionId
}

func RemoveSession(sessionId string) {
	delete(sessions, sessionId)
}

func sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
