package app2

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ServiceFuncHandler(w http.ResponseWriter, r *http.Request, app *JablkoApp) {
	pathVars := mux.Vars(r)
	log.Printf("%v\n", pathVars)
}
