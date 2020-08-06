package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func initSessionStore() {
	gob.Register(map[string]interface{}{})

	key := os.Getenv("RTL_STOREKEY")
	if key == "" {
		log.Fatal("session.init: RTL_STOREKEY not set in environment")
	}
	store = sessions.NewCookieStore([]byte(key))
}

func getSessionUsername(r *http.Request) string {
	session, _ := store.Get(r, "auth-session")

	profInterface, _ := session.Values["profile"]
	usernameInterface, _ := profInterface.(map[string]interface{})["name"]
	username, _ := usernameInterface.(string)

	log.Printf("%s is logged in", username)

	return username
}
