package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dspeirs7/animals/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (a *api) login(w http.ResponseWriter, r *http.Request) {
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

	session, err := a.store.Get(r, "session_token")
	session.Values["user"] = user.Username
	session.Save(r, w)

	if err != nil {
		a.errorResponse(w, r, http.StatusUnauthorized, err)
		return
	}

	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"sessionId": session.ID})
}

func (a *api) logout(w http.ResponseWriter, r *http.Request) {
	session, err := a.store.Get(r, "session_token")
	if err != nil {
		a.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	session.Options.MaxAge = -1
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
