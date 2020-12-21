// hanlders.go: MainApp's web handler and middleware definitions
// Cale Overstreet
// 2020/12/14
// Implementation of web routes for mainapp. Primary middleware 
// is the authenticationMiddleware.

package mainapp

import (
	"net/http"
	"context"
	"fmt"
	"log"
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
)

func (app *MainApp) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			// If path is login, send to login handler
			next.ServeHTTP(w, r)
			return 
		} else if r.URL.Path == "/logout" {
			next.ServeHTTP(w, r)
			return
		} else if strings.HasPrefix(r.URL.Path, "/local") {
			next.ServeHTTP(w, r)
			return
		}

		// Default values
		authenticated := false
		cookieValue := ""

		// First check if the key is present
		for _, val := range(r.Cookies()) {
			if val.Name == "jablkoLogin" {
				cookieValue = val.Value
				break;
			}
		}

		if cookieValue == "" {
			http.ServeFile(w, r, "./public_html/login/login.html")
			return
		}

		authenticated, sessionData, err := app.Db.ValidateSession(cookieValue)
		if err != nil {
			log.Println("ERROR: Unable to validate session.")
			log.Println(err)
		}

		if !authenticated {
			http.ServeFile(w, r, "./public_html/login/login.html")
			return
		}

		// How to pass data
		ctx := context.WithValue(r.Context(), "permissions", sessionData.Permissions) 
		ctx = context.WithValue(ctx, "username", sessionData.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *MainApp) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("User \"%s\" has requested dashboard (permission: %d)", r.Context().Value("username"), r.Context().Value("permissions"))

	// Read in dashboard template
	templateBytes, err := ioutil.ReadFile("./public_html/dashboard/template.html")
	if err != nil {
		log.Println("Unable to read template.html for dashboard")
	}

	template := string(templateBytes)

	// Read in toolbar
	toolbarBytes, err := ioutil.ReadFile("./public_html/toolbar/toolbar.html")
	if err != nil {
		log.Println("Unable to read template.html for dashboard")
		log.Println(err)
	}

	toolbar := string(toolbarBytes)

	var sb strings.Builder
	
	for _, modId := range app.ModHolder.Order {
		log.Println(modId)
		if curMod, ok := app.ModHolder.Mods[modId]; ok {
			sb.WriteString(curMod.Card(r))	
		}
	}

	replacer := strings.NewReplacer("$TOOLBAR", toolbar,
		"$JABLKO_MODULES", sb.String())

	w.Write([]byte(replacer.Replace(template)))
}

func (app *MainApp) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type loginHolder struct {
		Username string `json: "username"`
		Password string `json: "password"`
	}			

	var loginData loginHolder

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read login body")
		log.Println(err)
	}

	err = json.Unmarshal(body, &loginData)
	if err != nil {
		log.Println("Unable to unmarshal JSON data.")
		log.Println(err)
	}

	isCorrect, userData := app.Db.AuthenticateUser(loginData.Username, loginData.Password)

	if isCorrect {
		log.Println("User \"" + loginData.Username + "\" has logged in.")

		cookie, err := app.Db.CreateSession(loginData.Username, userData)
		if err != nil {
			log.Println("ERROR: Unable to create session for login")
			log.Println(err)
		}

		http.SetCookie(w, &cookie)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, `{"status": "good", "message": "Login succesful"}`)
	} else {
		w.Write([]byte(`{"status": "fail", "message": "Login data is wrong"}`))	
	}
}

func (app *MainApp) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookieValue := ""

	// First check if the key is present
	for key, val := range(r.Cookies()) {
		log.Println(key, val)

		if val.Name == "jablkoLogin" {
			cookieValue = val.Value
			break;
		}
	}

	if cookieValue == "" {
		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, `{"status": "fail", "message": "No matching cookie."}`)	
		return
	}

	err := app.Db.DeleteSession(cookieValue)	
	if err != nil {
		w.Header().Set("content-type", "application/json")
		fmt.Fprintln(w, `{"status": "fail", "message": "Failed to delete session."}`)	
		return
	}

	w.Header().Set("content-type", "application/json")
	fmt.Fprintln(w, `{"status": "good", "message": "Logged out."}`)	
}

func (app *MainApp) ModuleHandler(w http.ResponseWriter, r *http.Request) {
	// mod, func
	pathParams := mux.Vars(r)

	app.ModHolder.Mods[pathParams["mod"]].WebHandler(w, r)
}

func (app *MainApp) PublicHTMLHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	http.ServeFile(w, r, "./public_html/" + pathParams["pubdir"] + "/" + pathParams["file"])
}
