package app2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ccoverstreet/Jablko/core/database"
	"github.com/ccoverstreet/Jablko/core/modmanager"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type JablkoApp struct {
	server        *http.Server
	router        *mux.Router
	interruptChan chan os.Signal
	HTTPPort      int                       `json:"httpPort"`
	MessagingMods []string                  `json:"messagingMods"`
	ModM          *modmanager.ModManager    `json:"jmods"`
	DB            *database.DatabaseHandler `json:"-"`
}

func WrapRoute(route func(http.ResponseWriter, *http.Request, *JablkoApp), inst *JablkoApp) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		route(writer, req, inst)
	}
}

func WrapMiddleware(middleware func(http.Handler, *JablkoApp) http.Handler, app *JablkoApp) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware(next, app)
	}
}

// Creates an instance of JablkoApp with an initialized router
func CreateJablkoApp(config []byte) (*JablkoApp, error) {
	app := &JablkoApp{
		nil,
		nil,
		make(chan os.Signal, 2),
		80,
		[]string{},
		modmanager.NewModManager(),
		database.CreateDatabaseHandler(),
	}

	err := json.Unmarshal(config, &app)
	if err != nil {
		return nil, err
	}

	err = app.DB.LoadDatabase("./data/database.json")
	if err != nil {
		if err.Error() == "File does not exist" {
			log.Warn().
				Err(err).
				Msg("Unable to load existing database. Defaulting to empty database")

			// Initialize Empty Database
			app.DB.InitEmptyDatabase()
		} else {
			return nil, err
		}
	}

	log.Info().Msg("Loaded Database")
	app.router = mux.NewRouter()

	app.router.Use(WrapMiddleware(authMiddleware, app))
	app.router.Use(WrapMiddleware(loggingMiddleware, app))

	app.router.HandleFunc("/", WrapRoute(dashboardHandler, app))
	app.router.HandleFunc("/login", WrapRoute(userLoginHandler, app)).Methods("POST")
	app.router.HandleFunc("/admin", WrapRoute(adminPageHandler, app))
	app.router.HandleFunc("/admin/{func}", WrapRoute(AdminFuncHandler, app))
	app.router.HandleFunc("/assets/{file}", assetsHandler)
	app.router.PathPrefix("/jmod/").
		Handler(http.HandlerFunc(WrapRoute(passReqToJMOD, app))).
		Methods("GET", "POST")
	app.router.HandleFunc("/service/{func}", WrapRoute(ServiceFuncHandler, app))

	app.server = &http.Server{
		Addr:    ":" + strconv.Itoa(app.HTTPPort),
		Handler: app.router,
	}

	return app, nil
}

// Starts all JMODs that are loaded in JablkoApp.ModM
func (app *JablkoApp) StartJMODs() {
	app.ModM.StartAllJMODs()
}

// Starts JablkoApp HTTP server on port from loaded config
func (app *JablkoApp) Run() {
	go app.cleanup()
	signal.Notify(app.interruptChan, syscall.SIGINT, syscall.SIGTERM)
	app.server.ListenAndServe()
	app.interruptChan <- syscall.SIGINT
}

func (app *JablkoApp) cleanup() {
	<-app.interruptChan
	log.Info().
		Msg("Running cleanup code")
	log.Printf("%v\n", app.ModM.ProcMap)

	for jmodName, proc := range app.ModM.ProcMap {
		log.Info().
			Str("jmodName", jmodName).
			Msg("Stopping JMOD")

		proc.Stop()
	}

	app.server.Close()
}

func (app *JablkoApp) SaveConfig() error {
	log.Info().Msg("Saving Jablko Config...")

	config, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("./jablkoconfig.json", config, 0666)

	return err
}

// General HTTP error handler that makes handling code more concise
func handleHTTPError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%v", err)
}

// Logs request remote address, URI, and time taken to process.
// Time output is in microseconds
func loggingMiddleware(next http.Handler, app *JablkoApp) http.Handler {
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

// Authenticates inbound requests or validates that
// the requested route is open access
func authMiddleware(next http.Handler, app *JablkoApp) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Possible authentication clients are:
		// - Users
		// - JMODs
		// Which means that we should check for JMOD-KEY, JMOD-PORT headers or
		// jablko-session cookies

		// ---------- Allow safe routes through without auth ---------
		// Allow assets to be obtained without authentication
		if strings.HasPrefix(r.URL.String(), "/assets") {
			next.ServeHTTP(w, r)
			return
		}

		// Allow login requests to go through
		if r.URL.String() == "/login" { //strings.HasPrefix(r.URL.String(), "/login") {
			next.ServeHTTP(w, r)
			return
		}
		// --------- END safe routes ---------

		jmodKeyVal := r.Header.Get("JMOD-KEY")
		jmodPortVal, portValErr := strconv.Atoi(r.Header.Get("JMOD-PORT"))
		jablkoSessionCookie, cookieErr := r.Cookie("jablko-session")
		if jmodKeyVal == "" && cookieErr != nil {
			log.Warn().
				Msg("Unauthenticated request")

			if r.Method != "GET" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Not authenticated")
				return
			}

			loginPageHandler(w, r)
			return
		}

		// TODO: Refactor this shit
		if jmodKeyVal != "" {
			if portValErr != nil {
				log.Warn().
					Err(portValErr).
					Msg("Invalid port value")
				return
			}
			isValid, jmodName := app.ModM.IsValidService(jmodPortVal, jmodKeyVal)
			if !isValid {
				log.Warn().
					Msg("Unauthenticated request")

				fmt.Fprintf(w, "Not authenticated")
				return
			}

			r.Header.Set("JMOD-NAME", jmodName)
			r.Header.Set("Jablko-User-Permissions", strconv.Itoa(0))
		} else {
			isValid, permissionLevel := app.DB.ValidateSession(jablkoSessionCookie.Value)
			if !isValid {
				log.Warn().
					Msg("Session cookie not valid")

				if r.Method != "GET" {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Session invalid")
					return
				}

				loginPageHandler(w, r)
				return
			}

			r.Header.Set("Jablko-User-Permissions", strconv.Itoa(permissionLevel))
		}

		next.ServeHTTP(w, r)
	})
}

// Loads login page from file and sends it to client
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
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

// TODO: Should LoginUserHandler be a part of the database
// I am feeling no
func userLoginHandler(w http.ResponseWriter, r *http.Request, app *JablkoApp) {
	app.DB.LoginUserHandler(w, r)
}

// Generates dashboard HTML from all JMODs
func dashboardHandler(w http.ResponseWriter, r *http.Request, app *JablkoApp) {
	bTaskbar, err := ioutil.ReadFile("./html/taskbar.html")
	if err != nil {
		handleHTTPError(err, w)
		return
	}

	dashboardTemplate, err := ioutil.ReadFile("./core/app/template.html")
	if err != nil {
		handleHTTPError(err, w)
		return
	}

	WebCompStr, InstConfStr := app.ModM.GenerateJMODDashComponents()

	dashboardReplacer := strings.NewReplacer(
		"$JABLKO_TASKBAR", string(bTaskbar),
		"$JABLKO_WEB_COMPONENT_MAP_DEF", WebCompStr,
		"$JABLKO_JMOD_INSTANCE_CONF_MAP_DEF", InstConfStr,
	)

	fmt.Fprintf(w, "%s", dashboardReplacer.Replace(string(dashboardTemplate)))

	return
}

// Sends admin HTML page to client
func adminPageHandler(w http.ResponseWriter, r *http.Request, app *JablkoApp) {
	pageBytes, err := ioutil.ReadFile("./html/admin.html")
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

	fmt.Fprintf(w, "%s", strings.Replace(string(pageBytes), "$JABLKO_TASKBAR", string(bTask), 1))
}

// Acts as a reverse proxy to JMODs if the "/jmod" prefix is used on the route
func passReqToJMOD(w http.ResponseWriter, r *http.Request, app *JablkoApp) {
	err := app.ModM.PassRequest(w, r)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to pass request")
	}
}

var assetsMap = map[string]func(http.ResponseWriter) error{
	"standard.css": sendStandardCSS,
	"admin.js":     sendAdminJS,
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)

	assetName, ok := routeVars["file"]
	if !ok {
		handleHTTPError(fmt.Errorf("File not specified in route"), w)
		return
	}

	assetHandle, ok := assetsMap[assetName]
	if !ok {
		handleHTTPError(fmt.Errorf("No handler for requested asset"), w)
		return
	}

	err := assetHandle(w)
	if err != nil {
		handleHTTPError(err, w)
		return
	}
}

func sendStandardCSS(w http.ResponseWriter) error {
	fileBytes, err := ioutil.ReadFile("./html/standard.css")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/css")
	fmt.Fprintf(w, "%s", fileBytes)

	return nil
}

func sendAdminJS(w http.ResponseWriter) error {
	b, err := ioutil.ReadFile("./html/admin.js")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/javascript")
	fmt.Fprintf(w, "%s", b)

	return nil
}
