

package main

import (
	"net"
	"net/http"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"sync"
	"strings"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type testInstance struct {
	Id string `json:"id"`
	Source string `json:"source"`
	Value int `json:"value"`
}

type testConfig struct {
	sync.Mutex
	PortUDP int `json:"portUDP"`
	Instances []testInstance `json:"instances"`
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
	router := mux.NewRouter()
	router.HandleFunc("/webComponent", webComponentHandler)
	router.HandleFunc("/instanceData", instanceDataHandler)
	router.HandleFunc("/jmod/socket", SocketHandler)
	router.HandleFunc("/jmod/getUDPState", UDPStateHandler)
	router.HandleFunc("/jmod/{func}", JMODHandler).Methods("GET")

	port := os.Getenv("JABLKO_MOD_PORT")

	err := json.Unmarshal([]byte(os.Getenv("JABLKO_MOD_CONFIG")), &curConfig)
	if err != nil {
		panic(err)
	}
	log.Println(curConfig)

	// Start UDP server with in separate go routine
	// This server just prints the output and echoes
	go UDPServer()

	log.Println("Starting HTTP server...")
	http.ListenAndServe(":" + port, router)
}

func webComponentHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./webcomponent.js")
	if err != nil {
		fmt.Fprintf(w, "Unable to read WebComponent file")
	}

	fmt.Fprintf(w, "%s", b)
}

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

func JMODHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)

	sentBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	log.Println(string(sentBody))

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"hello": "From Tester"}`)
}

func UDPStateHandler(w http.ResponseWriter, r *http.Request) {
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

func UDPServer() {
	log.Println("Starting UDP Server...")
	serverAddr, err := net.ResolveUDPAddr("udp", ":49152")
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
