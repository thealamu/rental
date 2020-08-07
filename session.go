package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

var (
	errStoreFailure       = fmt.Errorf("Session store failure")
	errProfileNotFound    = fmt.Errorf("Profile not found")
	errSessionKeyNotFound = fmt.Errorf("Session key not found")
)

func initSessionStore() {
	gob.Register(map[string]interface{}{})

	key := os.Getenv("RTL_STOREKEY")
	if key == "" {
		log.Fatal("session.init: RTL_STOREKEY not set in environment")
	}
	store = sessions.NewCookieStore([]byte(key))
}

func invalidateSession(r *http.Request, w http.ResponseWriter) error {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		return errStoreFailure
	}

	session.Options.MaxAge = -1

	return session.Save(r, w)
}

func getProfileValue(r *http.Request, key string) (string, error) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		return "", errStoreFailure
	}

	profInterface, ok := session.Values["profile"]
	if !ok {
		return "", errProfileNotFound
	}
	rawValue, ok := profInterface.(map[string]interface{})[key]
	if !ok {
		return "", errSessionKeyNotFound
	}
	value, ok := rawValue.(string)
	if !ok {
		return "", errSessionKeyNotFound
	}

	return value, nil
}
