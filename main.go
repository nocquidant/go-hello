package main

import (
	"github.com/nocquidant/go-hello/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	// Version the version set by Makefile
	version string
	// BuildDate the build time set by Makefile
	buildDate string
)

// The single entry point
func main() {
	log.Printf("Build var 'version' is: %s", version)
	log.Printf("Build var 'time' is: %s", buildDate)
	cmd.Execute()
}
