package main

import (
	"net/http"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag := "middleware.auth"

		_, err := getProfileValue(r, "email")
		if err != nil {
			if err == errStoreFailure {
				respondError(tag, w, failCodeSessionDB, err.Error(), http.StatusInternalServerError)
			} else if err == errProfileNotFound {
				respondError(tag, w, failCodeAuth, err.Error(), http.StatusUnauthorized)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
