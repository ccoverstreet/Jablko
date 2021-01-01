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
	"fmt"
	"sync"
	"net/http"
	"strconv"
	"os"

	"github.com/gorilla/mux"

	"github.com/ccoverstreet/Jablko/src/mainapp"
	"github.com/ccoverstreet/Jablko/src/middleware"
	"github.com/ccoverstreet/Jablko/src/jlog"
)

const startingStatement = `Jablko Smart Home
Cale Overstreet
Version 0.3.0
License: GPLv3

`

func main() {
	fmt.Printf(startingStatement)

	initializeDirectories()

	// Create an instance of MainApp
	jablkoApp, err := mainapp.CreateMainApp("./jablkoconfig.json")
	if err != nil {
		jlog.Errorf("%v\n", err)
		jlog.Panic("Unable to create main app.")
	}

	router := initializeRoutes(jablkoApp)

	// TESTING SECTION
	for i := 0; i < 3; i++ {
		err = jablkoApp.ModHolder.InstallMod("builtin/interfacestatus")
		if err != nil {
			jlog.Errorf("Unable to install jablkomod.\n")
			jlog.Errorf("%v\n", err)
		}
	}

	err = jablkoApp.ModHolder.InstallMod("github.com/ccoverstreet/hamstermonitor-master")
	if err != nil {
		jlog.Errorf("%v\n", err)
	}

	// Start HTTP and HTTPS depending on Config
	// Wait for all to exit
	var wg sync.WaitGroup
	startJablko(jablkoApp, router, &wg)
	wg.Wait()
}

func initializeDirectories() {
	// Create tmp directory
	jlog.Printf("Making \"tmp\" directory...\n")
	err := os.Mkdir("./tmp", 0755)
	if err != nil {
		jlog.Warnf("Unable to make \"tmp\" directory: %v\n", err)
	}

	jlog.Printf("Checking if \"data\" dir exists...\n")
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		jlog.Printf("\"data\" directory not found. Creating directory.\n")
		err := os.Mkdir("./data", 0755)
		if err != nil {
			jlog.Errorf("Unable to make data directory: %v\n", err)
			jlog.Errorf("%v\n", err)
			panic(err)
		}
	}
}

func initializeRoutes(app *mainapp.MainApp) *mux.Router {
	jlog.Printf("Initializing routes...\n")
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
			jlog.Printf("Starting HTTP Server on Port %d\n", app.Config.HttpPort)	

			jlog.Printf("%v\n", http.ListenAndServe(":" + strconv.Itoa(app.Config.HttpPort), router))
		}()
	}

	if app.Config.HttpsPort > 1 {
		// Start https server
		wg.Add(1)
		go func() {
			defer wg.Done()
			jlog.Printf("Starting HTTPS Server on Port %d\n", app.Config.HttpsPort)	
		}()
	}

	return errs
}
