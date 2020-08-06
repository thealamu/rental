package main

import "net/http"

//check if user is loggedin already
func alreadyLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag := "middleware.alreadyloggedin"
		_, err := getProfileValue(r, "name")
		if err == nil {
			respondError(tag, w, failCodeBadParameter, "Already logged in", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
