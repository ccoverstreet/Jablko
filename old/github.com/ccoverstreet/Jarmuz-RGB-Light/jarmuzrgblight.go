// Jarmuz RGB Light
// Cale Overstreet
// May 10, 2021

/*
A Jablko Mod that communicates using UDP to RGB Lights. This
module is to serve as a demo for the full development chain
of Jablko.
*/

package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ccoverstreet/Jarmuz-RGB-Light/jablkodev"
	"github.com/gorilla/websocket"
)

const defaultConfig = `{
	"instances": [
		{
			"lightIPs": [
				"10.0.0.5"
			]
		}	
	]
}
`

type jmodConfig struct {
	sync.RWMutex
	Instances []instanceData `json:"instances"`
}

type instanceData struct {
	LightIPs []string `json:"lightIPs"`
}

// -------------------- GLOBALS --------------------
var globalConfig jmodConfig
var globalJMODKey string
var globalJMODPort string
var globalJablkoCorePort string

//go:embed webcomponent.js
var webcomponentFile []byte

// -------------------- END GLOBALS --------------------

func main() {
	// Get passed jmodKey. Used for authenticating jmods with Jablko
	globalJMODKey = os.Getenv("JABLKO_MOD_KEY")
	globalJablkoCorePort = os.Getenv("JABLKO_CORE_PORT")
	globalJMODPort = os.Getenv("JABLKO_MOD_PORT")

	// Get Passed config daata
	initConfig()
	log.Println(globalConfig)

	// Handles called by Jablko
	http.HandleFunc("/webComponent", WebComponentHandler)
	http.HandleFunc("/instanceData", InstanceDataHandler)
	http.HandleFunc("/jmod/socket", SocketHandler)

	log.Println(http.ListenAndServe(":"+globalJMODPort, nil))
}

func initConfig() {
	confStr := os.Getenv("JABLKO_MOD_CONFIG")
	log.Printf("\"%s\"", confStr)

	// Check if config was provided. Replace confStr with default
	// if not.
	if len(confStr) < 3 {
		log.Println("No config provided. Starting with default config")
		loadDefaultConfig()
		// Should also send a request to Jablko with updated config
		return
	}

	err := json.Unmarshal([]byte(confStr), &globalConfig)
	if err != nil {
		log.Printf("Provided config is invalid. Loading default config: %v", err)
	}
}

func loadDefaultConfig() {
	err := json.Unmarshal([]byte(defaultConfig), &globalConfig)
	if err != nil {
		log.Printf("FATAL ERROR: Default config is invalid")
		panic(err)
	}

	saveConfig()
}

// This function sends a JSON of the current config to Jablko
// which then triggers a config save on Jablko.
func saveConfig() {
	configBytes, err := json.Marshal(globalConfig)
	if err != nil {
		log.Printf("Unable to marshal config: %v", err)
		return
	}

	err = jablkodev.JablkoSaveConfig(globalJablkoCorePort, globalJMODPort, globalJMODKey, configBytes)
	if err != nil {
		log.Printf("ERROR: Unable to save config - %v", err)
	}
}

func WebComponentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", webcomponentFile)
}

func InstanceDataHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(globalConfig.Instances)
	if err != nil {
		log.Printf("Unable to generate JSON from globalConfig.Instances: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to generate JSON string for instances: %v", err)
		return
	}
	fmt.Fprintf(w, `%s`, b)
}

// WebSocketHandler
var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

// Main dashboard event handler
// Parses WebSocket data into condensed UDP values
// and sends packets to the target light
func SocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket handler called")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: Unable to upgrade WebSocket - %v", err)
	}
	defer conn.Close()

	// Populate map
	// SHOULD BE ABLE TO POPULATE MAP THROUGH WEBSOCKET
	// Maybe by sending requests of length 2
	connMap := make(map[string]*net.UDPConn)

	log.Println("Websocket connection established")
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ERROR: Error reading WebSocket message - %v", err)
			conn.WriteMessage(messageType, []byte(err.Error()))
			return
		}

		splitMessage := strings.Split(string(message), ",")

		if len(splitMessage) != 5 {
			log.Println("ERROR: Message is not of length 5")
			conn.WriteMessage(messageType, []byte("Message is not of length 5"))
			continue
		}

		rawAddr := splitMessage[0] + ":4123"

		// Check if connection already exists
		// If not, resolve the address and cache it
		if _, ok := connMap[rawAddr]; !ok {
			resAddr, err := net.ResolveUDPAddr("udp", rawAddr)
			if err != nil {
				log.Printf("ERROR: Unable to resolve UDP address of light - %v", err)
				conn.WriteMessage(messageType, []byte(err.Error()))
				return
			}

			light, err := net.DialUDP("udp", nil, resAddr)
			if err != nil {
				log.Printf("ERROR: Unable to dial UDP address - %v", err)
				conn.WriteMessage(messageType, []byte(err.Error()))
				return
			}

			connMap[rawAddr] = light
		}

		outBuf := [4]byte{0}

		// Write to light
		for i := 1; i < 5; i++ {
			val, err := strconv.Atoi(splitMessage[i])
			if err != nil {
				log.Printf("ERROR: Unable to convert to int - %v", err)
				conn.WriteMessage(messageType, []byte(err.Error()))
				continue
			}

			outBuf[i-1] = byte(val)
		}

		_, err = connMap[rawAddr].Write(outBuf[:])

		if err != nil {
			log.Printf("ERROR: Unable to write to light - %v", err)
			conn.WriteMessage(messageType, []byte(err.Error()))
			continue
		}

	}
}
