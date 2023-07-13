package middleware

import (
	"net/http"

	"github.com/dspeirs7/animals/internal/domain"
)

func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			sessions := domain.GetSessions()
			sessionId := cookie.Value
			session, ok := sessions[sessionId]
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			if session.IsExpired() {
				delete(sessions, sessionId)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
