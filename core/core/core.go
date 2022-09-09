package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ccoverstreet/Jablko/core/procmanager"
	"github.com/gorilla/mux"
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

	core.router.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile("./html/dashboard.html")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "%s", b)
	})

	core.router.HandleFunc("/admin/{func}", WrapRoute(AdminFuncHandler, core))
	core.router.HandleFunc("/assets/{file}", assetsHandler).Methods("GET")

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
	"admin.js":      {"admin.js", "text/javascript"},
	"standard.css":  {"standard.css", "text/css"},
	"dashboard.css": {"dashboard.css", "text/css"},
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]

	asset, ok := assetMap[file]
	if !ok {
		httpErrorHandler(w,
			CreateHTTPError(400, fmt.Sprintf("Invalid file '%s' requested", file), nil))
		return
	}

	b, err := os.ReadFile("html/" + asset.Path)
	if err != nil {
		httpErrorHandler(w,
			CreateHTTPError(500, fmt.Sprintf("Unable to serve file %s", file), nil))
		return
	}

	w.Header().Set("Content-Type", asset.MimeType)
	w.Write(b)
}
