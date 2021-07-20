// Jablko Core App
// Cale Overstreet
// Mar. 30, 2021

// Describes how the functionality of Jablko integrate
// into a single struct that is created in the main
// function.

package app

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/ccoverstreet/Jablko/core/database"
	"github.com/ccoverstreet/Jablko/core/modmanager"
)

type JablkoCoreApp struct {
	Router    *mux.Router               `json:"-"`
	HTTPPort  int                       `json:"httpPort"`
	ModM      *modmanager.ModManager    `json:"jmods"`
	DBHandler *database.DatabaseHandler `json:"-"`
}

func CreateJablkoCoreApp() *JablkoCoreApp {
	app := new(JablkoCoreApp)
	app.Router = mux.NewRouter()
	app.HTTPPort = 8080
	app.ModM = modmanager.NewModManager()
	app.DBHandler = database.CreateDatabaseHandler()

	app.Router.Use(app.LoggingMiddleware)
	app.Router.Use(app.AuthMiddleware)
	app.Router.HandleFunc("/", app.DashboardHandler).Methods("GET")
	app.Router.HandleFunc("/login", app.LoginHandler).Methods("POST")
	app.Router.HandleFunc("/logout", app.LogoutHandler).Methods("GET", "POST")
	app.Router.HandleFunc("/admin", app.AdminPageHandler).Methods("GET", "POST")
	app.Router.HandleFunc("/admin/{func}", app.AdminFuncHandler).Methods("GET", "POST")
	app.Router.HandleFunc("/service/{func}", app.ServiceHandler).Methods("GET", "POST")
	app.Router.HandleFunc("/jmod/{func}", app.PassToJMOD).Methods("GET", "POST")
	app.Router.HandleFunc("/assets/{file}", app.AssetsHandler).Methods("GET")

	return app
}

func (app *JablkoCoreApp) Init() error {
	// Runs through procedures to instantiate
	// config data.

	// Create data folder
	// Is a fatal error if this fails
	err := os.MkdirAll("./data", 0755)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to create data directory")
		panic(err)
	}

	log.Info().Msg("Loading Database...")
	err = app.DBHandler.LoadDatabase("./data/database.json")
	if err != nil {
		if err.Error() == "File does not exist" {
			log.Warn().
				Err(err).
				Msg("Unable to load existing database. Defaulting to empty database")

			// Initialize Empty Database
			app.DBHandler.InitEmptyDatabase()
		} else {
			panic(err)
		}
	}
	log.Info().Msg("Loaded Database")

	jablkoConfig, err := os.ReadFile("./jablkoconfig.json")
	if err != nil {
		log.Printf("%v\n", err)
		panic(err)
	}

	err = json.Unmarshal(jablkoConfig, app)
	if err != nil {
		log.Printf("%v\n", err)
		panic(err)
	}

	log.Printf("%v\n", app.ModM)
	log.Printf("%v\n", app)
	return app.ModM.StartAllJMODs()
}

func (app *JablkoCoreApp) SaveConfig() error {
	log.Info().Msg("Saving Jablko Config...")

	config, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("./jablkoconfig.json", config, 0666)

	log.Printf("%s\b", config)

	return nil
}

func (app *JablkoCoreApp) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		log.Info().
			Str("reqIPAddress", r.RemoteAddr).
			Str("URI", r.URL.String()).
			Int64("timeForReq", time.Since(start).Microseconds()).
			Msg("Logging Middleware")
	})
}

// Checks for jablko-session cookie
// Handles which route
func (app *JablkoCoreApp) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// Pass through to modmanager routes from jmods
		// This route is what JMODs use to request functions
		// from Jablko Core.
		if strings.HasPrefix(r.URL.String(), "/service") {
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

		isValid, permissionLevel := app.DBHandler.ValidateSession(jablkoSession.Value)

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

		r.Header.Set("Jablko-User-Permissions", strconv.Itoa(permissionLevel))
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

func (app *JablkoCoreApp) AdminPageHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./html/admin.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to read admin.html.")
		return
	}

	bTask, err := ioutil.ReadFile("./html/taskbar.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to read admin.html.")
		return
	}

	fmt.Fprintf(w, "%s", strings.Replace(string(b), "$JABLKO_TASKBAR", string(bTask), 1))
}

//go:embed template.html
var dashboardTemplate string

/// Creates the dashboard from the HTML template "index.html"
func (app *JablkoCoreApp) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	bTaskbar, err := ioutil.ReadFile("./html/taskbar.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	WebCompStr, InstConfStr := app.ModM.GenerateJMODDashComponents()

	dashboardReplacer := strings.NewReplacer(
		"$JABLKO_TASKBAR", string(bTaskbar),
		"$JABLKO_WEB_COMPONENT_MAP_DEF", WebCompStr,
		"$JABLKO_JMOD_INSTANCE_CONF_MAP_DEF", InstConfStr,
	)

	fmt.Fprintf(w, "%s", dashboardReplacer.Replace(dashboardTemplate))

	return
}

func (app *JablkoCoreApp) PassToJMOD(w http.ResponseWriter, r *http.Request) {
	err := app.ModM.PassRequest(w, r)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to pass request")
	}
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
	case "admin.js":
		b, err := ioutil.ReadFile("./html/admin.js")
		if err != nil {
			flagFail = true
			break
		}

		w.Header().Set("Content-Type", "text/javascript")
		fmt.Fprintf(w, "%s", b)
	default:
		flagFail = true
	}

	if flagFail {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Asset not found")
	}
}
