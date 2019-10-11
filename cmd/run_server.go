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
		msg  = fs.String("msg", "", "a custom message to display")
		_    = fs.String("config", "", "config file (optional)")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.JSONParser),
		ff.WithEnvVarPrefix("HELLO"))

	log.Printf("Flag var 'name' is: %s", *name)
	log.Printf("Flag var 'port' is: %d", *port)
	log.Printf("Flag var 'url' is: %s", *url)
	log.Printf("Flag var 'msg' is: %s", *msg)

	params := server.NewParameters(*name, *port, *url)
	if *msg != "" {
		params.WithMessage(*msg)
	}
	server.ListenAndServe(params)
}
