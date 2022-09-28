package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ccoverstreet/Jablko/core/procmanager"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type JablkoCore struct {
	PMan     procmanager.ProcManager `json:"mods"`
	PortHTTP int                     `json:"portHTTP"`
	router   *mux.Router
}

func CreateJablkoCore(config []byte) (*JablkoCore, error) {
	core := &JablkoCore{procmanager.CreateProcManager(), 8080, mux.NewRouter()}

	err := json.Unmarshal(config, core)
	if err != nil {
		return nil, err
	}

	// Add middleware
	core.router.Use(LoggingMiddleware)

	core.router.HandleFunc("/dashboard", wrapRoute(dashboardHandler))

	core.router.HandleFunc("/admin/{func}", WrapRoute(AdminFuncHandler, core))
	core.router.HandleFunc("/assets/{file}", assetsHandler).Methods("GET")
	core.router.PathPrefix("/mod/").
		Handler(http.HandlerFunc(WrapRoute(passRequestHandler, core))).
		Methods("GET", "POST")

	return core, nil
}

func createRouter(core *JablkoCore) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ASDASDASDASDASDASDAS")
	})

	return r
}

func WrapRoute(fun func(http.ResponseWriter, *http.Request, *JablkoCore), app *JablkoCore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fun(w, r, app)
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("route", r.URL.RequestURI()).
			Str("remoteAddr", r.RemoteAddr).
			Msg("Inbound request")

		next.ServeHTTP(w, r)
	})

}

func (core *JablkoCore) SaveConfig() error {
	b, err := json.Marshal(core)
	if err != nil {
		return err
	}

	return os.WriteFile("jablkoconfig.json", b, 0666)
}

func (core *JablkoCore) Start() {
	http.ListenAndServe(":"+strconv.Itoa(core.PortHTTP), core.router)
}

// Map of asset name to corresponding file
// Path is relative to html folder
type assetInfo struct {
	Path     string
	MimeType string
}

var assetMap = map[string]assetInfo{
	"jutil.js":      {"jutil.js", "text/javascript"},
	"admin.js":      {"admin.js", "text/javascript"},
	"standard.css":  {"standard.css", "text/css"},
	"dashboard.css": {"dashboard.css", "text/css"},
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]

	asset, ok := assetMap[file]
	if !ok {
		HTTPErrorHandler(w,
			CreateHTTPError(400, fmt.Sprintf("Invalid file '%s' requested", file), nil))
		return
	}

	b, err := os.ReadFile("html/" + asset.Path)
	if err != nil {
		HTTPErrorHandler(w,
			CreateHTTPError(500, fmt.Sprintf("Unable to serve file %s", file), nil))
		return
	}

	w.Header().Set("Content-Type", asset.MimeType)
	w.Header().Set("Cache-Control", "max-age=3600")
	w.Write(b)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request, core *JablkoCore) {
	core.PMan.GenerateWCScript()

	b, err := os.ReadFile("./html/dashboard.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

func passRequestHandler(w http.ResponseWriter, r *http.Request, core *JablkoCore) {
	params := r.URL.Query()
	modName, ok := params["mod"]
	if !ok {
		HTTPErrorHandler(w, CreateHTTPError(400,
			fmt.Sprintf("Mod not specified"), nil))
		return
	}

	err := core.PMan.PassRequest(modName[0], w, r)
	if err != nil {
		HTTPErrorHandler(w, CreateHTTPError(400,
			fmt.Sprintf("Unable to pass request"), err))
		return
	}

	// JSONResponse(w, struct{}{})
}
