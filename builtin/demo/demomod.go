package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const gDefaultConfig string = `
{
	"portUDP": 41111,
	"instances": [
		{
			"value": 32
		}
	]
}
`

// Frontend instances
type testInstance struct {
	Value int `json:"value"`
}

type testGlobalConfig struct {
	sync.Mutex
	PortUDP   int            `json:"portUDP"`
	Instances []testInstance `json:"instances"`
}

type stateUDP struct {
	sync.Mutex
	Data string
}

// Global data
var curConfig testGlobalConfig
var curStateUDP = stateUDP{sync.Mutex{}, "Startup Message"}
var jablkoCorePort string
var jablkoModPort string
var jablkoModKey string

//go:embed webcomponent.js
var webcomponentFile []byte

func main() {
	// Grabbing settings from environment variable
	jablkoCorePort = os.Getenv("JABLKO_CORE_PORT")
	jablkoModPort = os.Getenv("JABLKO_MOD_PORT")
	jablkoModKey = os.Getenv("JABLKO_MOD_KEY")

	// Pull in config from environment variable and marshal into
	// state struct. Load default config and save if value is not
	// present
	configStr := os.Getenv("JABLKO_MOD_CONFIG")
	loadConfig(configStr)

	router := mux.NewRouter()

	// Required Routes
	router.HandleFunc("/webComponent", webComponentHandler) // Route called by Jablko
	router.HandleFunc("/instanceData", instanceDataHandler) // Route called by Jablko

	// Application Routes
	router.HandleFunc("/jmod/socket", SocketHandler)        // Application route for WebSockets
	router.HandleFunc("/jmod/getUDPState", UDPStateHandler) // Simple GET for UDP Server State
	router.HandleFunc("/jmod/testConfigSave", TestConfigSave)

	log.Println(curConfig.PortUDP)
	// Start UDP server with in separate go routine
	// This server just prints the output and echoes
	go UDPServer(curConfig.PortUDP)

	log.Println("Starting HTTP server...")
	log.Println(http.ListenAndServe(":"+jablkoModPort, router))
}

func loadConfig(config string) {
	// If no config is provided
	if len(config) < 3 {
		err := json.Unmarshal([]byte(gDefaultConfig), &curConfig)

		log.Println(curConfig)

		if err != nil {
			panic(err)
		}

		err = JablkoSaveConfig(jablkoCorePort, jablkoModPort, jablkoModKey, []byte(gDefaultConfig))
		if err != nil {
			panic(err)
		}

		return
	}

	err := json.Unmarshal([]byte(config), &curConfig)

	if err != nil {
		panic(err)
	}
}

// The webcomponent handler returns the javascript for a WebComponent
// javascript class. In production, the file should be cached so that
// disk reads are kept to a minimum
func webComponentHandler(w http.ResponseWriter, r *http.Request) {
	// Example for debugging, reads file on every request
	// Leave this commented out
	b, err := ioutil.ReadFile("./webcomponent.js")
	if err != nil {
		fmt.Fprintf(w, "Unable to read WebComponent file")
	}
	fmt.Fprintf(w, "%s", b)

	//fmt.Fprintf(w, "%s", webcomponentFile)
}

// Instance data returns a javascript object string with
// keys representing individual instances and sub objects
// representing instance data
func instanceDataHandler(w http.ResponseWriter, r *http.Request) {
	curConfig.Lock()
	defer curConfig.Unlock()

	log.Println(curConfig)
	log.Println(curConfig.Instances)
	data, err := json.Marshal(curConfig.Instances)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to marshal config data")
		return
	}

	fmt.Fprintf(w, "%s", data)
}

// Returns the most recent UDP message received by the UDP Server
// to the Jablko client
func UDPStateHandler(w http.ResponseWriter, r *http.Request) {
	// Shows how to get user permissions from r header
	fmt.Println("User Permission Level:", r.Header.Get("Jablko-User-Permissions"))
	curStateUDP.Lock()
	defer curStateUDP.Unlock()
	fmt.Fprintf(w, `{"state": "%s"}`, curStateUDP.Data)
}

// ---------- WEB SOCKETS ----------
// Example for implementation of Web Sockets
// The CheckOrigin method of the upgrader
// must be ignored to as the origin of the
// request is modified by the Jablko Core
// proxy
var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("User Permission Level:", r.Header.Get("Jablko-User-Permissions"))
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Received: %s\n", message)
		response := append(message, []byte(" received by server")...)
		err = conn.WriteMessage(messageType, response)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// ---------- UDP Server ----------
// UDP communication is often better suited for
// IoT applications. Each Jablko Mod can start
// a UDP server/client to communicate with pmods.
type restartFlag struct {
	sync.Mutex
	Restart bool
}

func UDPServer(port int) {
	log.Println("Starting UDP Server...")
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
		}
		x := string(buf[0:n])
		curStateUDP.Lock()
		curStateUDP.Data = strings.Replace(string(buf[0:n]), "\n", "", -1)
		curStateUDP.Unlock()

		// Echo data
		serverConn.WriteToUDP([]byte("ECHO: "+x), addr)

		log.Println("From Client:", string(buf[0:n]))
	}
}

// Marshals current state into JSON string and sends
// to Jablko. Jablko then saves the updated data to the config
// file
func TestConfigSave(w http.ResponseWriter, r *http.Request) {
	curConfig.Lock()
	defer curConfig.Unlock()

	b, err := json.Marshal(curConfig)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Couldn't marshal state to string")
		return
	}

	err = JablkoSaveConfig(jablkoCorePort, jablkoModPort, jablkoModKey, b)

	log.Println(err)
}
