

package main

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/jmod/stateless/{modId}/data", GETDataHandler).Methods("GET")
	router.HandleFunc("/jmod/stateless/{modId}/socket", SocketHandler)
	router.HandleFunc("/jmod/{state}/{modId}/{modRoute}", JModHandler)

	fmt.Printf("\nTESTER: %s\n\n", os.Environ())

	port := os.Getenv("JABLKO_MOD_PORT")

	http.ListenAndServe(":" + port, router)
}

// 
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
		}

		fmt.Printf("Received: %s\n", message)
		err = conn.WriteMessage(messageType, []byte("Received by server"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

