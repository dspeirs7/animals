package middleware

import "net/http"

func CommonMiddleware(h http.Handler) http.Handler {
	return Logger(Session(h))
}
