package domain

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"
)

var sessions = make(map[string]Session)

type Session struct {
	Username string
	Expiry   time.Time
}

func GetSession(id string) (Session, bool) {
	session, ok := sessions[id]
	return session, ok
}

func SetSession(session Session) string {
	sessionId := sessionId()
	sessions[sessionId] = session
	return sessionId
}

func RemoveSession(sessionId string) {
	delete(sessions, sessionId)
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
