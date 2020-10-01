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
	"io/ioutil"
	"sync"
	"log"
	"net/http"
	"strconv"
	"github.com/buger/jsonparser"
)

const startingStatement = `Jablko Smart Home
Cale Overstreet
Version 2.0.0
License: GPLv3

`

type jablkoConfig struct {
	httpPort int
	httpsPort int
}

var config = jablkoConfig{httpPort: 8080, httpsPort: -1}

func main() {
	log.Printf(startingStatement)

	initializeConfig()
	log.Printf("%v\n", config)

	initializeRoutes()

	// Start HTTP and HTTPS depending on config
	// Wait for all to exit
	var wg sync.WaitGroup
	startJablko(config, &wg)
	wg.Wait()
}

func initializeConfig() {
	configData, err := ioutil.ReadFile("./jablko_config.json")
	if err != nil {
		log.Printf("%v\n", err)
		panic("Error opening and reading config file\n")
	}

	// Get HTTP data
	httpPort, err := jsonparser.GetInt(configData, "http", "port")
	if err != nil {
		log.Printf("%v\n", err)
		panic("Error getting HTTP port data\n")
	}

	config.httpPort = int(httpPort)

	httpsPort, err := jsonparser.GetInt(configData, "https", "port")
	if err != nil {
		log.Printf("HTTPS port config not set in config file\n")
	} else {
		config.httpsPort = int(httpsPort)
	}
}

func initializeRoutes() {
	http.HandleFunc("/", dashboardHandler)
}

func startJablko(config jablkoConfig, wg *sync.WaitGroup) chan error {
	errs := make(chan error)

	if config.httpPort > 1 {
		// Start http port
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting HTTP Server on Port %d\n", config.httpPort)	

			log.Printf("%v\n", http.ListenAndServe(":" + strconv.Itoa(config.httpPort), nil))
		}()
	}

	if config.httpsPort > 1 {
		// Start https server
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting HTTPS Server on Port %d\n", config.httpsPort)	
		}()
	}

	return errs
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "./public_html/dashboard/dashboard_template.html")
	}
}
