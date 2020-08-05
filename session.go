package main

import (
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	key := os.Getenv("RTL_STOREKEY")
	store = sessions.NewCookieStore([]byte(key))
}
