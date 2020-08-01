package main

import (
	"os"

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

	srv := newServer(port)
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
