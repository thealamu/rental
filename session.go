package main

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	key := os.Getenv("RTL_STOREKEY")
	if key == "" {
		log.Fatal("session.init: RTL_STOREKEY not set in environment")
	}
	store = sessions.NewCookieStore([]byte(key))
}
