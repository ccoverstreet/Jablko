

package main

import (
	"net/http"
	"fmt"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/jmod/{state}/{modId}/{modRoute}", JModHandler)

	fmt.Printf("\nTESTER: %s\n\n", os.Environ())

	port := os.Getenv("JABLKO_MOD_PORT")

	http.ListenAndServe(":" + port, router)
}

func JModHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	fmt.Fprintf(w, `{"hello": "From Tester"}`)
}
