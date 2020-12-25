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
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/ccoverstreet/Jablko/src/jlog"
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
			jlog.Errorf("ERROR: Unable to validate session.\n")
			jlog.Errorf("%v\n", err)
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
	jlog.Printf("User \"%s\" has requested dashboard (permission: %d)\n", r.Context().Value("username"), r.Context().Value("permissions"))

	// Read in dashboard template
	templateBytes, err := ioutil.ReadFile("./public_html/dashboard/template.html")
	if err != nil {
		jlog.Warnf("Unable to read template.html for dashboard\n")
		jlog.Warnf("%v\n", err)
	}

	template := string(templateBytes)

	// Read in toolbar
	toolbarBytes, err := ioutil.ReadFile("./public_html/toolbar/toolbar.html")
	if err != nil {
		jlog.Warnf("Unable to read template.html for dashboard\n")
		jlog.Warnf("%v\n", err)
	}

	toolbar := string(toolbarBytes)

	var sb strings.Builder
	
	for _, modId := range app.ModHolder.Order {
		jlog.Println(modId)
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
		jlog.Warnf("Unable to read login body\n")
		jlog.Warnf("%v\n")
	}

	err = json.Unmarshal(body, &loginData)
	if err != nil {
		jlog.Warnf("Unable to unmarshal JSON data.\n")
		jlog.Println("%v\n", err)
	}

	isCorrect, userData := app.Db.AuthenticateUser(loginData.Username, loginData.Password)

	if isCorrect {
		jlog.Println("User \"" + loginData.Username + "\" has logged in.\n")

		cookie, err := app.Db.CreateSession(loginData.Username, userData)
		if err != nil {
			jlog.Errorf("ERROR: Unable to create session for login\n")
			jlog.Errorf("%v\n", err)
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
	for _, val := range(r.Cookies()) {
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
