package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/schicho/mensa-restful/internal"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

// whichAddress returns the address for the server to listen on.
// The address is either taken from the PORT environment variable
// or from the -port flag.
// The environment variable takes precedence.
// If neither is set, the default port 8080 is used.
func whichAddress(args []string) (string, error) {
	flags := flag.NewFlagSet("mensa-restful", flag.ContinueOnError)
	var (
		port = flags.String("port", "8080", "port to listen on")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return "", err
	}

	var addr string

	if envPort := os.Getenv("PORT"); envPort != "" {
		addr = "0.0.0.0:" + envPort
	} else {
		addr = "0.0.0.0:" + *port
	}
	return addr, nil
}

func run(args []string) error {
	addr, err := whichAddress(args)
	if err != nil {
		return err
	}

	srv, err := internal.NewServer()
	if err != nil {
		return err
	}

	log.Printf("listening on %v\n", addr)
	return http.ListenAndServe(addr, srv)
}
