package main

import (
	"github.com/nocquidant/go-hello/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	// Version the version set by Makefile
	Version string
	// BuildTime the build time set by Makefile
	BuildTime string
)

// The single entry point
func main() {
	log.Printf("Build var 'version' is: %s", Version)
	log.Printf("Build var 'time' is: %s", BuildTime)
	cmd.Execute()
}
