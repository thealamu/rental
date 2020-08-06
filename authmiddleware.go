package main

import (
	"log"
	"net/http"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag := "middleware.auth"

		session, err := store.Get(r, "auth-session")
		if err != nil {
			log.Printf("%s: %v", tag, err)
			respondError(tag, w, failCodeSessionDB, "Session store failure", http.StatusInternalServerError)
			return
		}

		_, ok := session.Values["profile"]
		if !ok {
			respondError(tag, w, failCodeAuth, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
