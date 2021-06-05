// Jablko Core
// Cale Overstreet
// Mar. 30, 2021

// Core Jablko process entrypoint. Responsible for
// proxying requests, authentication, process
// management

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ccoverstreet/Jablko/core/app"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	fmt.Printf("%s\b", startingStatement)

	jablkoApp := new(app.JablkoCoreApp)
	err := jablkoApp.Init()
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Error in initialization")
	}

	// Setup process cleanup handler
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Info().
			Msg("Cleaning up child processes")

		jablkoApp.ModM.CleanProcesses()
		os.Exit(0)
	}()

	log.Info().Msg("Starting HTTP Server")
	log.Fatal().Err(http.ListenAndServe(":8080", jablkoApp.Router))
}

const startingStatement = `
Jablko Smart Home
Cale Overstreet
Version 0.3.0
License: GPLv3

`
