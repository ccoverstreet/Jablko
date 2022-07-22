package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func AdminRouteHandler(w http.ResponseWriter, r *http.Request, core *JablkoCore) {
	fun := mux.Vars(r)["func"]

	var res []byte
	var err error

	switch fun {
	case "stopJMOD":
		res, err = stopJMOD(r, core)
	default:
		log.Println("Implement")
	}

	log.Println(res, err)
}

func stopJMOD(r *http.Request, core *JablkoCore) ([]byte, error) {

	return nil, nil
}
