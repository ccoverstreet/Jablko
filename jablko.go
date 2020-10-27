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
	"time"
	"github.com/gorilla/mux"
	"github.com/buger/jsonparser"

	"github.com/ccoverstreet/Jablko/jablkomods"
)

const startingStatement = `Jablko Smart Home
Cale Overstreet
Version 2.0.0
License: GPLv3

`

type jablkoConfig struct {
	httpPort int
	httpsPort int
	moduleOrder []string
}

var config = jablkoConfig{httpPort: 8080, httpsPort: -1}

func main() {
	log.Printf(startingStatement)

	initializeConfig()
	log.Printf("%v\n", config)

	router := initializeRoutes()

	// Start HTTP and HTTPS depending on config
	// Wait for all to exit
	var wg sync.WaitGroup
	startJablko(config, router, &wg)
	wg.Wait()
}

func initializeConfig() {
	configData, err := ioutil.ReadFile("./jablkoconfig.json")
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

	jablkoModulesSlice, _, _, err := jsonparser.Get(configData, "jablkoModules")
	if err != nil {
		panic("Error get Jablko Modules Config\n")
	}

	err = jablkomods.Initialize(jablkoModulesSlice)
	if err != nil {
		log.Println("Error initializing Jablko Mods")
		log.Println(err)
	}

	// Initialize module order in config file
	moduleOrderSlice, _, _, err := jsonparser.Get(configData, "moduleOrder")

	jsonparser.ArrayEach(moduleOrderSlice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		config.moduleOrder = append(config.moduleOrder, string(value))
	})

	log.Println(config)

	log.Println(jablkomods.ModMap)
}

func initializeRoutes() *mux.Router {
	r := mux.NewRouter()

	// Timing Middleware
	r.Use(timingMiddleware)

	r.HandleFunc("/", dashboardHandler).Methods("GET")

	return r
}

func startJablko(config jablkoConfig, router *mux.Router, wg *sync.WaitGroup) chan error {
	errs := make(chan error)

	if config.httpPort > 1 {
		// Start http port
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting HTTP Server on Port %d\n", config.httpPort)	

			log.Printf("%v\n", http.ListenAndServe(":" + strconv.Itoa(config.httpPort), router))
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

func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Calling next handler
		next.ServeHTTP(w, r)

		end := time.Now()

		log.Printf("Request took %d ns\n", end.Sub(start))
	})
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public_html/dashboard/dashboard_template.html")

	var x http.Request
	
	for i := 0; i < len(config.moduleOrder); i++ {
		log.Println(jablkomods.ModMap[config.moduleOrder[i]].Card(&x))	
	}
}
