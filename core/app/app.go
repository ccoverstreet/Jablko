// Jablko Core App
// Cale Overstreet
// Mar. 30, 2021

// Describes how the functionality of Jablko integrate
// into a single struct that is created in the main
// function.

package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/ccoverstreet/Jablko/core/database"
	"github.com/ccoverstreet/Jablko/core/modmanager"
)

type JablkoCoreApp struct {
	Router    *mux.Router
	ModM      *modmanager.ModManager
	DBHandler *database.DatabaseHandler
}

func (app *JablkoCoreApp) Init() error {
	// Runs through procedures to instantiate
	// config data.
	app.initRouter()

	// Read jablkoconfig.json
	confByte, err := ioutil.ReadFile("./jablkoconfig.json")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to read jablkoconfig.json")

		return err
	}

	sourceConf, _, _, err := jsonparser.Get(confByte, "sources")
	if err != nil {
		panic(err)
	}

	log.Info().Msg("Creating module manager...")
	newModM, err := modmanager.NewModManager(sourceConf)
	if err != nil {
		panic(err)
	}
	app.ModM = newModM
	log.Info().Msg("Created module manager")

	log.Info().Msg("Creating database handler...")
	app.DBHandler = database.CreateDatabaseHandler()
	err = app.DBHandler.LoadDatabase("./data/database.json")
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Unable to load existing database")
	}
	log.Info().Msg("Created database handler")

	return nil
}

func (app *JablkoCoreApp) initRouter() {
	// Creates the gorilla/mux router passed to
	// http.ListenAndServe

	router := mux.NewRouter()
	router.Use(app.LoggingMiddleware)
	router.Use(app.AuthMiddleware)
	router.HandleFunc("/", app.DashboardHandler).Methods("GET")
	router.HandleFunc("/login", app.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", app.LogoutHandler).Methods("GET", "POST")
	router.HandleFunc("/jmod/{func}", app.PassToJMOD).Methods("GET", "POST")
	router.HandleFunc("/assets/{file}", app.AssetsHandler).Methods("GET")

	app.Router = router
}

func (app *JablkoCoreApp) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("reqIPAddress", r.RemoteAddr).
			Str("URI", r.URL.String()).
			Msg("Logging Middleware")

		next.ServeHTTP(w, r)
	})
}

// Checks for jablko-session cookie
func (app *JablkoCoreApp) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("Authentication Middleware")
		jablkoSession, err := r.Cookie("jablko-session")

		// Allow assets to be obtained without authentication
		if strings.HasPrefix(r.URL.String(), "/assets") {
			next.ServeHTTP(w, r)
			return
		}

		// Allow login requests to go through
		if strings.HasPrefix(r.URL.String(), "/login") {
			next.ServeHTTP(w, r)
			return
		}

		// Unable to get cookie
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Session cookie not found")

			if r.Method != "GET" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Not logged in")
				return
			}

			app.LoginPageHandler(w, r)

			return
		}

		isValid := app.DBHandler.ValidateSession(jablkoSession.Value)

		if !isValid {
			log.Warn().
				Msg("Session cookie not valid")

			if r.Method != "GET" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Session invalid")
				return
			}

			app.LoginPageHandler(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *JablkoCoreApp) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./html/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Cannot read login.html")
		log.Error().
			Err(err).
			Msg("Unable to read login.html")
	}

	fmt.Fprintf(w, "%s", b)
}

func (app *JablkoCoreApp) LoginHandler(w http.ResponseWriter, r *http.Request) {
	app.DBHandler.LoginUserHandler(w, r)
}

func (app *JablkoCoreApp) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jablko-session")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid cookie value")
	}

	app.DBHandler.DeleteSession(cookie.Value)
}

func (app *JablkoCoreApp) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./html/index.html")
	if err != nil {
		return
	}

	template := string(b)

	builderWC := strings.Builder{}
	builderInstance := strings.Builder{}

	for modSource, subProc := range app.ModM.ProcMap {
		fmt.Printf("%s: %v\n", modSource, subProc)
		b1, err := app.getWebComponent(subProc.Port)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Unable to get WebComponent")
			continue
		}
		builderWC.WriteString("\njablkoWebCompMap[\"" + modSource + "\"] = ")
		builderWC.Write(b1)

		b2, err := app.getInstanceData(subProc.Port)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Unable to get JMOD instance data")
			continue
		}
		builderInstance.WriteString("\njablkoInstanceConfMap[\"" + modSource + "\"] = ")
		builderInstance.Write(b2)
	}

	dashboardReplacer := strings.NewReplacer(
		"$JABLKO_WEB_COMPONENT_MAP_DEF", builderWC.String(),
		"$JABLKO_JMOD_INSTANCE_CONF_MAP_DEF", builderInstance.String(),
	)

	fmt.Fprintf(w, "%s", dashboardReplacer.Replace(template))

	return
}

func (app *JablkoCoreApp) PassToJMOD(w http.ResponseWriter, r *http.Request) {
	// Checks for JMOD_Source URL parameter
	// Returns 404
	source := r.FormValue("JMOD-Source")
	log.Info().
		Str("JMOD", source).
		Str("URI", r.URL.RequestURI()).
		Msg("Passing request to JMOD")

	// Check if no JMOD-Source header value found
	if len(source) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Empty JMOD-Source parameter")
		log.Warn().
			Str("JMOD-Source", source).
			Msg("Empty JMOD-Source parameter")
		return
	}

	// Check if JMOD-Source is a valid option
	if _, ok := app.ModM.ProcMap[source]; !ok {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "JMOD not found")
		return
	}

	app.ModM.PassRequest(w, r)
	return

}

func (app *JablkoCoreApp) getWebComponent(modPort int) ([]byte, error) {
	resp, err := http.Get("http://localhost:" + strconv.Itoa(modPort) + "/webComponent")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (app *JablkoCoreApp) getInstanceData(modPort int) ([]byte, error) {
	resp, err := http.Get("http://localhost:" + strconv.Itoa(modPort) + "/instanceData")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (app *JablkoCoreApp) AssetsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	flagFail := false

	switch vars["file"] {
	case "standard.css":
		b, err := ioutil.ReadFile("./html/standard.css")
		if err != nil {
			flagFail = true
			break
		}

		w.Header().Set("Content-Type", "text/css")
		fmt.Fprintf(w, "%s", b)
	default:
		flagFail = true
	}

	if flagFail {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Asset not found")
	}
}
