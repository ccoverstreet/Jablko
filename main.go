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

	"github.com/ccoverstreet/Jablko/core/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	fmt.Printf("%s\b", startingStatement)

	// Make Directories if don't exist
	err := os.MkdirAll("./log", 0700)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to make log directory")
	}

	err = os.MkdirAll("./data", 0700)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to make log directory")
	}

	err = os.MkdirAll("./tmp", 0700)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to make tmp directory")
	}

	jablkoApp := app.CreateJablkoCoreApp()
	err = jablkoApp.LoadConfig()
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Error in loading config")

		panic(err)
	}

	// Start JMODs
	errs := jablkoApp.StartJMODs()
	if len(errs) != 0 {
		log.Error().Msg("Error when starting JMODs. See prior logs.")
	}

	fmt.Println(jablkoApp)

	log.Info().Msg("Starting HTTP Server")
	log.Error().
		Err(http.ListenAndServe(":8080", jablkoApp.Router)).
		Msg("Jablko stopping")
}

const startingStatement = `
Jablko Smart Home
Cale Overstreet
Version 0.3.0
License: GPLv3

`
