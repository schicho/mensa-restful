package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/schicho/mensa-restful/internal"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		port = flags.Int("port", 8080, "port to listen on")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	srv, err := internal.NewServer()
	if err != nil {
		return err
	}
	fmt.Printf("listening on :%d\n", *port)
	return http.ListenAndServe(addr, srv)
}

