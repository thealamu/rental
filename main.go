package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var port string

func main() {
	parseFlags()
	if port == "" {
		//Check if port is set in the environment
		if port = os.Getenv("PORT"); port == "" {
			// It's not, default to 8080
			port = "8080"
		}
	}
	fmt.Println(port)
}

func parseFlags() {
	flag.StringVarP(&port, "port", "p", "", "port to start the server on")
	flag.Parse()
}
