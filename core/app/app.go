package app

import (
	"encoding/json"
	"net/http"

	"github.com/ccoverstreet/Jablko/core/modmanager"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type JablkoCore struct {
	ModM   modmanager.ModManager `json:"jmods"`
	router *mux.Router
}

func WrapRoute(route func(http.ResponseWriter, *http.Request, *JablkoCore), core *JablkoCore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		route(w, r, core)
	}
}

func CreateHTTPRouter(core *JablkoCore) *mux.Router {
	r := &mux.Router{}
	r.HandleFunc("/", WrapRoute(dashboardHandler, core))

	r.PathPrefix("/jmod/").
		Handler(http.HandlerFunc(WrapRoute(PassReqToJMOD, core))).
		Methods("GET", "POST")

	return r
}

func CreateJablkoCore(config []byte) (*JablkoCore, error) {
	newApp := &JablkoCore{}
	err := json.Unmarshal(config, newApp)

	log.Printf("%v", newApp)

	newApp.router = CreateHTTPRouter(newApp)

	return newApp, err
}

func (core *JablkoCore) StartAllMods() {
	core.ModM.StartAll()
}

func (core *JablkoCore) Listen() {
	log.Info().Msg("Jablko Core online and listening.")
	http.ListenAndServe(":8080", core.router)
}

func (core *JablkoCore) Cleanup() {
	log.Info().Msg("Cleaning up Jablko Core processes.")
	err := core.ModM.Cleanup()

	if err != nil {
		log.Error().
			Err(err).
			Msg("Errors occured when cleaning up Jablko Core processes")
	}
}

func PassReqToJMOD(w http.ResponseWriter, r *http.Request, core *JablkoCore) {
	core.ModM.PassReqToJMOD(w, r)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request, core *JablkoCore) {
	core.ModM.GetDashboard()
}
