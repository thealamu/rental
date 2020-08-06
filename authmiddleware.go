package main

import (
	"log"
	"net/http"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag := "middleware.auth"

		_, err := getProfileValue(r, "email")
		if err != nil {
			if err == errStoreFailure {
				log.Printf("%s: %v", tag, err)
				respondError(tag, w, failCodeSessionDB, "Session store failure", http.StatusInternalServerError)
			} else if err == errProfileNotFound {
				log.Printf("%s: %v", tag, err)
				respondError(tag, w, failCodeAuth, "Unauthorized", http.StatusUnauthorized)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
