package cmd

import (
	"flag"
	"os"

	"github.com/nocquidant/go-hello/pkg/server"
	"github.com/peterbourgon/ff"
	log "github.com/sirupsen/logrus"
)

// Execute an HTTP server to serve webhooks as REST endpoints
func Execute() {
	fs := flag.NewFlagSet("go-hello", flag.ExitOnError)
	var (
		name = fs.String("name", "hello-svc", "name of the service")
		port = fs.Int("port", 8484, "http port to listent to")
		url  = fs.String("url", "localhost:8485/hello", "remote url to get")
	)

	ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("HELLO"))

	log.Printf("Flag var 'name' is: %s", *name)
	log.Printf("Flag var 'port' is: %d", *port)
	log.Printf("Flag var 'url' is: %s", *url)

	params := server.NewParameters(*name, *port, *url)
	server.ListenAndServe(params)
}
