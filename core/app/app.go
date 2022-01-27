package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ccoverstreet/Jablko/core/modmanager"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type JablkoCore struct {
	ModM   modmanager.ModManager `json:"jmods"`
	router *mux.Router
}

func CreateHTTPRouter() *mux.Router {
	r := &mux.Router{}
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ASDASDASD")
	})

	return r
}

func CreateJablkoCore(config []byte) (*JablkoCore, error) {
	newApp := &JablkoCore{}
	err := json.Unmarshal(config, newApp)

	log.Printf("%v", newApp)

	newApp.router = CreateHTTPRouter()

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
