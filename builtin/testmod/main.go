package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/mod/{func}", demoHandler)

	log.Println("Starting HTTP Server")
	http.ListenAndServe(":9090", router)
	fmt.Println("vim-go")
}

func demoHandler(w http.ResponseWriter, r *http.Request) {
	fun := mux.Vars(r)["func"]

	log.Printf("Function received: %s", fun)
	fmt.Fprintf(w, "{'msg': 'received'}")
}
