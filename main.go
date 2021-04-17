// Jablko Core
// Cale Overstreet
// Mar. 30, 2021

// Core Jablko process entrypoint. Responsible for
// proxying requests, authentication, process
// management

package main

import (
	"net/http"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ccoverstreet/Jablko/core/app"
)

func main() {
	fmt.Printf("%s\b", startingStatement)

	jablkoApp := new(app.JablkoCoreApp)
	err := jablkoApp.Init()
	if err != nil {
		log.Panic().
			Err(err).
			Caller().
			Msg("ASd")
	}

	err = jablkoApp.ModManager.StartJablkoMod("builtin/test")

	log.Fatal().Err(http.ListenAndServe(":8080", jablkoApp.Router))
}

const startingStatement = `
Jablko Smart Home
Cale Overstreet
Version 0.3.0
License: GPLv3

`
