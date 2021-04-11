// Jablko Core
// Cale Overstreet
// Mar. 30, 2021

// Core Jablko process entrypoint. Responsible for
// proxying requests, authentication, process
// management

package main

import (
	"log"
	"net/http"


	"github.com/ccoverstreet/Jablko/core/logging"
	"github.com/ccoverstreet/Jablko/core/app"
)



var jablkoApp *app.JablkoCoreApp

func main() {
	setupLogging()

	log.Println(startingStatement)

	jablkoApp = new(app.JablkoCoreApp)
	err := jablkoApp.Init()
	if err != nil {
		panic(err)
	}

	err  = jablkoApp.ModManager.StartJablkoMod("./builtin/test")
	log.Println(err)

	log.Println(jablkoApp.ModManager.SubprocessMap)


	log.Fatal(http.ListenAndServe(":8080", jablkoApp.Router))
}

func setupLogging() {
	log.SetFlags(0)
	log.SetOutput(new(logging.JablkoLogger))
}


const startingStatement = `
Jablko Smart Home
Cale Overstreet
Version 0.3.0
License: GPLv3

`
