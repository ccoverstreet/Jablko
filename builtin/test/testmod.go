

package main

import (
	"net"
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/fart", homeHandler)
	router.HandleFunc("/jmod/socket", SocketHandler)
	router.HandleFunc("/jmod/{func}", GETDataHandler).Methods("GET")
	//router.HandleFunc("/jmod/{state}/{modId}/{modRoute}", JModHandler)

	fmt.Printf("\nTESTER: %s\n\n", os.Environ())

	port := os.Getenv("JABLKO_MOD_PORT")

	// Start UDP server with in separate go routine
	// This server just prints the output and echoes
	go UDPServer()

	fmt.Println("Starting HTTP")
	http.ListenAndServe(":" + port, router)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ASDASDASDASDASD FROM TESTMOD")
}

func JModHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)

	sentBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(sentBody))

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"hello": "From Tester"}`)
}

func GETDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"reqType": "GET", "res": "Test Module GET Response"}`)
}

// ---------- WEB SOCKETS ----------
// Example for implementation of Web Sockets
// The CheckOrigin method of the upgrader 
// must be ignored to as the origin of the
// request is modified by the Jablko Core
// proxy
var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool {return true}}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SOCKET HANDLER")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Received: %s\n", message)
		err = conn.WriteMessage(messageType, []byte("Received by server"))
		if err != nil {
			fmt.Println(err)
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

	fmt.Println("ASD")
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		x := string(buf[0:n])

		// Echo data
		serverConn.WriteToUDP([]byte("ECHO: " + x), addr)

		fmt.Println("From Client:", string(buf[0:n]))
	}
}
