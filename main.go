// Jablko Core
// Cale Overstreet
// Mar. 30, 2021

// github.com/ccoverstreet/Jablko
// Core Jablko process entrypoint. Responsible for
// proxying requests, authentication, process
// management.

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ccoverstreet/Jablko/core/app2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	fmt.Printf("%s\b", startingStatement)

	// Initialize directory structure that Jablko uses
	CreateDirectories()

	configBytes, err := ioutil.ReadFile("./jablkoconfig.json")
	if err != nil {
		log.Warn().Msg("jablkoconfig.json not found. Starting Jablko with default config")
		configBytes = []byte(app2.JABLKO_DEFAULT_CONFIG)
		err = ioutil.WriteFile("jablkoconfig.json", configBytes, 0644)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Unable to save default config to file")
		}
	}

	jablkoApp2, err := app2.CreateJablkoApp(configBytes)
	if err != nil {
		panic(err)
	}
	log.Printf("%v\n", jablkoApp2)

	jablkoApp2.StartJMODs()
	jablkoApp2.Run()
}

func CreateDirectories() {
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
}

const startingStatement = `
Jablko Smart Home
Cale Overstreet
Version 0.3.0
License: GPLv3

`
