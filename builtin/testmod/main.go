package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/mod/{func}", demoHandler)
	router.HandleFunc("/webComponent", webComponentHandler)

	log.Println("Starting HTTP Server")
	http.ListenAndServe(":9090", router)
	fmt.Println("vim-go")
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func demoHandler(w http.ResponseWriter, r *http.Request) {
	fun := mux.Vars(r)["func"]

	res := struct {
		Msg  string `json:"msg"`
		Func string `json:"func"`
	}{
		"Test message received",
		fun,
	}

	log.Printf("Function %s requested", fun)

	sendJSONResponse(w, res)
}

func webComponentHandler(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile("webcomponent.js")
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"err": "Unable to read webcomponent file"}`)
		return
	}

	w.Header().Set("Content-Type", "text/javascript")
	w.Write(b)
}
