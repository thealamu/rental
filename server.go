package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

var errNoPort = fmt.Errorf("No port specified")

func newServer(port string, router http.Handler) *http.Server {
	if port == "" {
		log.Fatal("newServer: ", errNoPort)
	}

	return &http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: router,
	}
}

//shutdown handles gracefully shutting down a server
func shutdown(s *http.Server) {
	log.Printf("server.shutdown: shutting down server on %s\n", s.Addr)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server.shutdown: ", err)
	}
}

func run(s *http.Server) {
	log.Printf("server.run: starting server on %s\n", s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("server.run: ", err)
	}
}
