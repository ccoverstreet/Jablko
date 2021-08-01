package service

import (
	"log"
	"net/http"

	"github.com/ccoverstreet/Jablko/core/app2"
	"github.com/gorilla/mux"
)

func ServiceHandler(w http.ResponseWriter, r *http.Request, app *app2.JablkoApp) {
	pathVars := mux.Vars(r)
	log.Printf("%v\n", pathVars)
}
