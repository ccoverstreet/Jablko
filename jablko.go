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
	"strings"
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/buger/jsonparser"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ccoverstreet/Jablko/jablkomods"
	"github.com/ccoverstreet/Jablko/src/database"
)

const startingStatement = `Jablko Smart Home
Cale Overstreet
Version 0.3.0
License: GPLv3

`

type generalConfig struct {
	HttpPort int
	HttpsPort int
	ConfigMap map[string]string
	ModuleOrder []string
}

var jablkoConfig = generalConfig{HttpPort: 8080, HttpsPort: -1}

type MainApp struct {} // Placeholder struct for implementing the JablkoInterface interface

var Jablko MainApp // Creating the MainApp instance
 
func (jablko MainApp) Tester() {
	log.Println("Shit")
}

var jablkoDB *sql.DB // Database handle

func (jablko MainApp) SyncConfig(modId string) {
	log.Println("Initial")
	log.Println(jablkoConfig.ConfigMap)

	ConfigTemplate:= `{
	"http": {
		"port": $httpPort
	},
	"https": {
		"port": $httpsPort
	},
	"jablkoModules": {
		$moduleString
	},
	"moduleOrder": [
		$moduleOrder
	]
}
`

	log.Println(ConfigTemplate)

	if _, ok := jablkomods.ModMap[modId]; !ok {
		log.Printf("Cannot find module %s", modId)
		return 
	}

	testStr, err := jablkomods.ModMap[modId].ConfigStr()
	if err != nil {
		log.Printf("Unable to get Config string for module %s\n", modId)
	}


	jablkoConfig.ConfigMap[modId] = string(testStr)

	log.Println(string(testStr))

	log.Println("Updated")
	log.Println(jablkoConfig.ConfigMap)

	// Create JSON to dump to Config file

	// Prepare each module's string
	jablkoModulesStr := ""
	index := 0
	for key, value := range jablkoConfig.ConfigMap {
		if index > 0 {
			jablkoModulesStr = jablkoModulesStr + ",\n\t\t\"" + key + "\":" + value
		} else {
			jablkoModulesStr = jablkoModulesStr + "\"" + key + "\":" + value
		}

		index = index + 1
	}

	// Prepare Module Order
	orderStr := ""
	for index, val := range jablkoConfig.ModuleOrder {
		if index > 0 {
			orderStr = orderStr + ",\n" + "\t\t\"" + val + "\""
		} else {
			orderStr = orderStr + "\"" + val + "\""
		}
	}

	log.Println(orderStr)

	r := strings.NewReplacer("$httpPort", strconv.Itoa(jablkoConfig.HttpPort),
	"$httpsPort", strconv.Itoa(jablkoConfig.HttpsPort),
	"$moduleString", jablkoModulesStr,
	"$moduleOrder", orderStr)

	ConfigDumpStr := r.Replace(ConfigTemplate)

	err = ioutil.WriteFile("./jablkoconfig.json", []byte(ConfigDumpStr), 0022)
	if err != nil {
		log.Println(`Unable to write to "jablkoconfig.json".`)
	}
}

func main() {
	log.Printf(startingStatement)

	initializeConfig()

	router := initializeRoutes()

	jablkoDB = database.Initialize()
	defer jablkoDB.Close()

	// Start HTTP and HTTPS depending on Config
	// Wait for all to exit
	var wg sync.WaitGroup
	startJablko(jablkoConfig, router, &wg)
	wg.Wait()
}

func initializeConfig() {
	ConfigData, err := ioutil.ReadFile("./jablkoconfig.json")
	if err != nil {
		log.Printf("%v\n", err)
		panic("Error opening and reading Config file\n")
	}

	// Get HTTP data
	httpPort, err := jsonparser.GetInt(ConfigData, "http", "port")
	if err != nil {
		log.Printf("%v\n", err)
		panic("Error getting HTTP port data\n")
	}

	jablkoConfig.HttpPort = int(httpPort)

	httpsPort, err := jsonparser.GetInt(ConfigData, "https", "port")
	if err != nil {
		log.Printf("HTTPS port Config not set in Config file\n")
	} else {
		jablkoConfig.HttpsPort = int(httpsPort)
	}

	jablkoModulesSlice, _, _, err := jsonparser.Get(ConfigData, "jablkoModules")
	if err != nil {
		panic("Error get Jablko Modules Config\n")
	}

	ConfigMap, err := jablkomods.Initialize(jablkoModulesSlice, Jablko)
	if err != nil {
		log.Println("Error initializing Jablko Mods")
		log.Println(err)
	}

	jablkoConfig.ConfigMap = ConfigMap

	// Initialize module order in Config file
	moduleOrderSlice, _, _, err := jsonparser.Get(ConfigData, "moduleOrder")

	jsonparser.ArrayEach(moduleOrderSlice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		jablkoConfig.ModuleOrder = append(jablkoConfig.ModuleOrder, string(value))
	})

	// Print Config
	log.Println(jablkoConfig)

	// Print module map
	log.Println(jablkomods.ModMap)
}

func initializeRoutes() *mux.Router {
	r := mux.NewRouter()

	// Timing Middleware
	r.Use(timingMiddleware)
	r.Use(authenticationMiddleware)

	r.HandleFunc("/", dashboardHandler).Methods("GET")
	r.HandleFunc("/jablkomods/{mod}/{func}", moduleHandler).Methods("POST")

	return r
}

func initializeDatabase() {
	}

func startJablko(jablkoConfig generalConfig, router *mux.Router, wg *sync.WaitGroup) chan error {
	errs := make(chan error)

	if jablkoConfig.HttpPort > 1 {
		// Start http port
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting HTTP Server on Port %d\n", jablkoConfig.HttpPort)	

			log.Printf("%v\n", http.ListenAndServe(":" + strconv.Itoa(jablkoConfig.HttpPort), router))
		}()
	}

	if jablkoConfig.HttpsPort > 1 {
		// Start https server
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting HTTPS Server on Port %d\n", jablkoConfig.HttpsPort)	
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

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Println(r.Form)

		authenticated := false
		cookieValue := ""

		for key, val := range(r.Cookies()) {
			log.Println(key, val)

			if val.Name == "jablkologin" {
				cookieValue = val.Value
				break;
			}
		}

		if cookieValue == "" {
			log.Println("No login cookie found.")
		}

		log.Println(authenticated)

		/*
		if val := r.Cookies()[0] {
			log.Println(val)
		} else {
			log.Println("ASDASDASD")
		}
		*/
		
		next.ServeHTTP(w, r)
	})
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "./public_html/dashboard/dashboard_template.html")

	var x http.Request

	var sb strings.Builder
	
	for i := 0; i < len(jablkoConfig.ModuleOrder); i++ {
		sb.WriteString(jablkomods.ModMap[jablkoConfig.ModuleOrder[i]].Card(&x))	
	}

	cookie := http.Cookie {
		Name: "jablkologin",
		Value: "11111",
		Expires: time.Now().Add(1 * time.Minute),
	}

	http.SetCookie(w, &cookie)

	w.Write([]byte(sb.String()))
}

func moduleHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL.Path)
	splitPath := strings.Split(r.URL.Path, "/")
	if val, ok := jablkomods.ModMap[splitPath[2]]; ok {
		val.WebHandler(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid path received."}`))
	}
}
