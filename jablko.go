// jablko.go: Entrypoint for Jablko
// Cale Overstreet
// October 1, 2020
/* 
This is the entrypoint for the Jablko Smart Home system 
and start of the Go-based system. This project is switching 
to Go to greatly improve performance, code 
readability, and overall architecture. The restrictions of
the Go language (and features) will allow for a more
performant and usable Jablko. The chatbot functionality can 
also be greatly improved using goroutines and threading.
A bottleneck in the NodeJS version was any computationally 
intense route.
*/

package main

import (
	"sync"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ccoverstreet/Jablko/src/mainapp"
	"github.com/ccoverstreet/Jablko/src/middleware"
	"github.com/ccoverstreet/Jablko/src/jlog"
)

const startingStatement = `Jablko Smart Home
Cale Overstreet
Version 0.2.0
License: GPLv3

`

func main() {
	log.Printf(startingStatement)

	jlog.Printf("Hello\n")
	jlog.Warnf("Warning Log\n")
	jlog.Errorf("Error Log\n")

	// Create an instance of MainApp
	jablkoApp, err := mainapp.CreateMainApp("./jablkoconfig.json")
	if err != nil {
		log.Panic("Unable to create main app.")
	}

	router := initializeRoutes(jablkoApp)

	log.Println(jablkoApp)
	log.Println(jablkoApp.ModHolder)

	// Start HTTP and HTTPS depending on Config
	// Wait for all to exit
	var wg sync.WaitGroup
	startJablko(jablkoApp, router, &wg)
	wg.Wait()
}

func initializeRoutes(app *mainapp.MainApp) *mux.Router {
	r := mux.NewRouter()

	// Timing Middleware
	r.Use(middleware.TimingMiddleware)
	r.Use(app.AuthenticationMiddleware)

	r.HandleFunc("/", app.DashboardHandler).Methods("GET")

	r.HandleFunc("/jablkomods/{mod}/{func}", app.ModuleHandler).Methods("POST")
	r.HandleFunc("/local/{mod}/{func}", app.ModuleHandler).Methods("POST")
	r.HandleFunc("/{pubdir}/{file}", app.PublicHTMLHandler).Methods("GET")
	r.HandleFunc("/login", app.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", app.LogoutHandler).Methods("POST")

	return r
}

func startJablko(app *mainapp.MainApp, router *mux.Router, wg *sync.WaitGroup) chan error {
	errs := make(chan error)

	if app.Config.HttpPort > 1 {
		// Start http port
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting HTTP Server on Port %d\n", app.Config.HttpPort)	

			log.Printf("%v\n", http.ListenAndServe(":" + strconv.Itoa(app.Config.HttpPort), router))
		}()
	}

	if app.Config.HttpsPort > 1 {
		// Start https server
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting HTTPS Server on Port %d\n", app.Config.HttpsPort)	
		}()
	}

	return errs
}
