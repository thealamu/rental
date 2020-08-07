package main

import (
	"log"
	"net/http"
)

func handleLogout(w http.ResponseWriter, r *http.Request) {
	tag := "handler.logout"

	logoutURL, err := newLogoutURL(r)
	if err != nil {
		respondError(tag, w, failCodeUnknown, err.Error(), http.StatusInternalServerError)
		return
	}

	email, err := getProfileValue(r, "email")
	if err != nil {
		respondError(tag, w, failCodeSessionDB, err.Error(), http.StatusInternalServerError)
		return
	}

	//Invalidate the session
	err = invalidateSession(r, w)
	if err != nil {
		respondError(tag, w, failCodeSessionDB, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("%s: %s logged out", tag, email)

	http.Redirect(w, r, logoutURL.String(), http.StatusTemporaryRedirect)
}
