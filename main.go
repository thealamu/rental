package main

import (
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
		//Check if port is set in the environment
		if port = os.Getenv("PORT"); port == "" {
			//It's not, default to 8080
			port = "8080"
		}
	}

	ctx, cancel := signalcontext.OnInterrupt()
	defer cancel()

	//Register paths
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler)

	//Path /cars
	carsRouter := router.PathPrefix("/cars").Subrouter()
	carsRouter.HandleFunc("", getPublicCars).Methods(http.MethodGet)
	carsRouter.HandleFunc("/{car_id}", getSinglePublicCar).Methods(http.MethodGet)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	srv := newServer(port, loggedRouter)
	//Run server
	go run(srv)

	<-ctx.Done()
	//Gracefully shutdown
	shutdown(srv)
}

func parseFlags() {
	flag.StringVarP(&port, "port", "p", "", "port to start the server on")
	flag.Parse()
}
