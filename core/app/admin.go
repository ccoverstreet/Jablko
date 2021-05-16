package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (app *JablkoCoreApp) AdminFuncHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch vars["func"] {
	case "addUser":
		log.Printf("ASDASD")
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid admin function requested")
	}
}
