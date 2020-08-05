package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sethvargo/go-signalcontext"
	flag "github.com/spf13/pflag"
)

var port string

func main() {
	parseFlags()
	if port == "" {
		//If port is not set in the environment too, default to 8080
		if port = os.Getenv("PORT"); port == "" {
			port = "8080"
		}
	}

	//Create a context that closes when SIGINT is received
	ctx, cancel := signalcontext.OnInterrupt()
	defer cancel()

	router := mux.NewRouter()
	router.HandleFunc("/", getCommonEndpoints).Methods(http.MethodGet)

	//Path /cars
	carsRouter := router.PathPrefix("/cars").Subrouter()
	carsRouter.HandleFunc("", getPublicCars).Methods(http.MethodGet)
	carsRouter.HandleFunc("/{car_id:[0-9]+}", getSinglePublicCar).Methods(http.MethodGet)

	//Path /merchants/{merchant}
	router.HandleFunc("/merchants/{merchant:[a-zA-Z ]+}", getSingleMiniMerchant)

	//Path /auth/login
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", handleLogin)
	authRouter.HandleFunc("/login/callback", handleLoginCallback)

	//Path /merchants/me
	merchmeRouter := router.PathPrefix("/merchants/me").Subrouter()
	merchmeRouter.HandleFunc("", getMerchantMe)

	//Wrap the root router with one that logs every request
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	//Init the db, so we know of any errors before we start handling requests
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Fatalf("main.Main: %v during database init", err)
	}
	defer db.close()

	//prepare the session store
	initSessionStore()

	srv := newServer(port, loggingRouter)
	go run(srv)

	<-ctx.Done()
	//We received SIGINT, gracefully shutdown
	shutdown(srv)
}

func parseFlags() {
	flag.StringVarP(&port, "port", "p", "", "port to start the server on")
	flag.Parse()
}
