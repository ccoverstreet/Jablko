

package main

import (
	"net"
	"net/http"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"sync"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Frontend instances
type testInstance struct {
	Value int `json:"value"`
}

type testConfig struct {
	sync.Mutex
	PortUDP int `json:"portUDP"`
	Instances map[string]testInstance `json:"instances"`
}

// 
type stateUDP struct {
	sync.Mutex
	Data string
}

// Global instances
var curConfig testConfig
var curStateUDP = stateUDP{sync.Mutex{}, "Startup Message"}

func main() {
	curConfig.Instances = make(map[string]testInstance)

	router := mux.NewRouter()
	// Required Routes
	router.HandleFunc("/webComponent", webComponentHandler) // Route called by Jablko
	router.HandleFunc("/instanceData", instanceDataHandler) // Route called by Jablko

	// Application Routes
	router.HandleFunc("/jmod/socket", SocketHandler) // Application route for WebSockets
	router.HandleFunc("/jmod/getUDPState", UDPStateHandler) // Simple GET for UDP Server State


	// Pull in port for running HTTP server that communicates with Jablko
	port := os.Getenv("JABLKO_MOD_PORT")
	log.Println(port)

	// Pull in config from environment variable and marshal into
	// state struct
	err := json.Unmarshal([]byte(os.Getenv("JABLKO_MOD_CONFIG")), &curConfig)
	if err != nil {
		panic(err)
	}
	log.Println(curConfig)

	// Start UDP server with in separate go routine
	// This server just prints the output and echoes
	go UDPServer(curConfig.PortUDP)

	log.Println("Starting HTTP server...")
	http.ListenAndServe(":" + port, router)
}

// The webcomponent handler returns the javascript for a WebComponent
// javascript class. In production, the file should be cached so that
// disk reads are kept to a minimum
func webComponentHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./webcomponent.js")
	if err != nil {
		fmt.Fprintf(w, "Unable to read WebComponent file")
	}

	fmt.Fprintf(w, "%s", b)
}

// Instance data returns a javascript object string with
// keys representing individual instances and sub objects 
// representing instance data
func instanceDataHandler(w http.ResponseWriter, r *http.Request) {
	curConfig.Lock()
	defer curConfig.Unlock()

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
var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool {return true}}

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
	serverAddr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(port))
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
		serverConn.WriteToUDP([]byte("ECHO: " + x), addr)

		log.Println("From Client:", string(buf[0:n]))
	}
}
