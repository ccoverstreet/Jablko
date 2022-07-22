package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("JMOD example")

	r := &mux.Router{}
	r.HandleFunc("/webcomponent", WebComponentHandler)
	r.HandleFunc("/jmod/{func}", JMODRouteHandler)

	http.ListenAndServe(":8080", r)
}

func WebComponentHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./webcomponent.js")
	if err != nil {
		log.Printf("ERROR: Unable to read webcomponent file - %v", err)
	}

	fmt.Fprintf(w, "%s", b)
}

func JMODRouteHandler(w http.ResponseWriter, r *http.Request) {
	routeFunc := mux.Vars(r)["func"]

	log.Println(routeFunc)
}
