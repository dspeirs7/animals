package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dspeirs7/animals/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (a *api) login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		decoder := json.NewDecoder(r.Body)
		var user domain.User

		if err := decoder.Decode(&user); err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		admin, err := a.userRepo.GetUser(ctx, "admin")
		if err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(user.Password)); err != nil {
			a.errorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		expiresAt := time.Now().Add(3600 * time.Second)
		sessionId := domain.SetSession(domain.Session{Username: admin.Username, Expiry: expiresAt})

		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: sessionId, Expires: expiresAt, Path: "/", SameSite: http.SameSiteLaxMode})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"sessionId": sessionId})
	case http.MethodOptions:
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *api) logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cookie, err := r.Cookie("session_token")
		if err != nil {
			a.logger.Info(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sessionId := cookie.Value
		delete(sessions, sessionId)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	case http.MethodOptions:
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
