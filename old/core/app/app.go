package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

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

func WrapMiddleware(middleware func(http.Handler, *JablkoCore) http.Handler, core *JablkoCore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware(next, core)
	}
}

func LoggingMiddleware(next http.Handler, core *JablkoCore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Str("remoteAddress", r.RemoteAddr).
			Msg("Incoming request")

		next.ServeHTTP(w, r)
	})
}

func CreateHTTPRouter(core *JablkoCore) *mux.Router {
	r := &mux.Router{}

	r.Use(WrapMiddleware(LoggingMiddleware, core))
	r.Use(WrapMiddleware(AuthMiddleware, core))

	r.HandleFunc("/", WrapRoute(dashboardHandler, core))

	r.PathPrefix("/jmod/").
		Handler(http.HandlerFunc(WrapRoute(PassReqToJMOD, core))).
		Methods("GET", "POST")

	r.HandleFunc("/admin/{func}", WrapRoute(AdminRouteHandler, core))

	r.HandleFunc("/assets/{file}", WrapRoute(assetsHandler, core))

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

func routeErrorHandler(w http.ResponseWriter, err error) {
	_, filename, line, _ := runtime.Caller(1)

	log.Error().
		Str("file", filename).
		Int("line", line).
		Err(err).
		Msg("Error occured while handling route.")

	http.Error(w, err.Error(), 406)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request, core *JablkoCore) {
	html := core.ModM.GetDashboard()

	fmt.Fprintf(w, "%s", html)
}

var VALID_ASSETS = [...]string{
	"common.css",
	"home.svg",
}

var MIMEMAP = map[string]string{
	".css": "text/css",
	".js":  "text/javascript",
	".svg": "image/svg+xml",
}

func isValidAsset(filename string) bool {
	for _, a := range VALID_ASSETS {
		if a == filename {
			return true
		}
	}

	return false
}

func assetsHandler(w http.ResponseWriter, r *http.Request, core *JablkoCore) {
	filename := mux.Vars(r)["file"]

	if !isValidAsset(filename) {
		routeErrorHandler(w, fmt.Errorf("Invalid asset '%s' requested", filename))
		return
	}

	data, err := ioutil.ReadFile("html/" + filename)
	if err != nil {
		routeErrorHandler(w, err)
		return
	}

	ext := filepath.Ext(filename)

	mime := MIMEMAP[ext]
	w.Header().Add("Content-Type", MIMEMAP[ext])

	if strings.HasPrefix(mime, "text") {
		fmt.Fprintf(w, "%s", data)
	} else {
		fmt.Fprintf(w, "%s", data)
	}

}
