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
			if r.URL.Path == "/" {
				http.ServeFile(w, r, "./public_html/login/login.html")
			} else {
				w.Header().Set("Content-Type","application/json")	
				fmt.Fprintf(w, `{"status":"fail", "fart"}`)
			}

			return
		}

		authenticated, sessionData, err := app.Db.ValidateSession(cookieValue)
		if err != nil {
			jlog.Errorf("ERROR: Unable to validate session.\n")
			jlog.Errorf("%v\n", err)
		}

		if !authenticated {
			if r.URL.Path == "/" {
				http.ServeFile(w, r, "./public_html/login/login.html")
			} else {
				w.Header().Set("Content-Type", "application/json")	
				fmt.Fprintf(w, `{"status":"fail", "fart"}`)
			}

			return
		}

		// How to pass data
		ctx := context.WithValue(r.Context(), "permissions", sessionData.Permissions) 
		ctx = context.WithValue(ctx, "username", sessionData.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *MainApp) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Generates dashboard HTML page from javascript fragments
	// provided by Jablko Mods. 
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

	// String builders for source code
	var sourceMap strings.Builder
	sourceMap.WriteString("\nconst jablko_sourceMap = {\n")

	var moduleCardConfig strings.Builder
	moduleCardConfig.WriteString("\nconst jablko_configArr = [\n")

	sourceFileCache := make(map[string]bool) // Used to check if source has already been added
	isFirstMod := true // Used just for formatting
	
	// Loop through mods and gather source files
	for _, modId := range app.ModHolder.Order {
		modSourcePath := app.ModHolder.Mods[modId].SourcePath()

		if _, ok := sourceFileCache[modSourcePath]; !ok {
			if (!isFirstMod) {
				sourceMap.WriteString(",\n")
			}

			sourceBytes, err := ioutil.ReadFile("./" + modSourcePath + "/jablkomod.js")
			if err != nil {
				jlog.Warnf("Unable to read source \"%s\".\n", modSourcePath)
				jlog.Warnf("%v\n", err)
				continue
			}

			sourceMap.WriteString(`"` + modSourcePath + `":` + string(sourceBytes))
			
			sourceFileCache[modSourcePath] = true
		}

		if curMod, ok := app.ModHolder.Mods[modId]; ok {
			if (!isFirstMod) {
				moduleCardConfig.WriteString(",\n")
			}

			moduleCardConfig.WriteString(`{"mod": "` + modSourcePath + `" ,"config": ` + curMod.ModuleCardConfig() + `}`)	
		} else {
			jlog.Warnf("Dashboard card not available for \"%s\". Module not found.\n", modId)
		}

		isFirstMod = false
	}
	
	sourceMap.WriteString("\n};\n")

	moduleCardConfig.WriteString("\n];\n")

	// Replacer generates the final HTML
	replacer := strings.NewReplacer("$TOOLBAR", toolbar,
		"$JABLKO_SOURCE_MAP", sourceMap.String(), 
		"$JABLKO_INIT_ARR", moduleCardConfig.String()) 

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
		jlog.Warnf("%v\n", err)
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

	// Check if handler
	if curMod, ok := app.ModHolder.Mods[pathParams["mod"]]; ok {
		curMod.WebHandler(w, r)
	} else {
		jlog.Errorf("Module \"%s\" not found in module map. Please check if it is installed\n", pathParams["mod"])
	}
}

func (app *MainApp) PublicHTMLHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	http.ServeFile(w, r, "./public_html/" + pathParams["pubdir"] + "/" + pathParams["file"])
}

func (app *MainApp) AdminHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	
	ctx := r.Context()
	permissions, ok := ctx.Value("permissions").(int)
	if !ok {
		jlog.Warnf("Permissions field incorrect. Access denied.\n")
	} 

	if permissions < 2 {
		jlog.Warnf("User not authorized for this action. Ignoring request.\n")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"fail", "message":"Not authorized"}`)
		return 
	}

	if pathParams["func"] == "addMod" {
		addMod(app, w, r)
	} else if pathParams["func"] == "deleteMod" {
		deleteMod(app, w, r)	
	} else if pathParams["func"] == "addUser" {
		addUser(app, w, r)
	} else if pathParams["func"] == "deleteUser" {
	} else if pathParams["func"] == "updateMod" {
		updateMod(app, w, r)
	} else if pathParams["func"] == "getModConfig" {
		getModConfig(app, w, r)
	} else if pathParams["func"] == "registerMod" {
		jlog.Warnf("Registering mods has not been implemented\n")
	}
}

func addMod(app *MainApp, w http.ResponseWriter, r *http.Request) {
	type addModBody struct {
		SourcePath string `json:"sourcePath"`
	}

	w.Header().Set("Content-Type", "application/json")

	var parsedBody addModBody

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jlog.Errorf("Unable to read \"/admin/addMod\" body.\n")
		jlog.Errorf("%v\n", err)
		fmt.Fprintf(w, `{"status":"fail","message":"Unable to read request body."}`)
		return
	}

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		jlog.Warnf("Unable to unmarshal JSON data.\n")
		jlog.Println("%v\n", err)
		fmt.Fprintf(w, `{"status":"fail","message":"Unable to parse JSON body."}`)
		return
	}


	err = app.ModHolder.InstallMod(parsedBody.SourcePath)
	if err != nil {
		jlog.Errorf("%v\n", err)
		fmt.Fprintf(w, `{"status":"fail","message":"` + err.Error() + `"}`)
		return
	}

	fmt.Fprintf(w, `{"status": "good"}`)
}

func deleteMod(app *MainApp, w http.ResponseWriter, r *http.Request) {
	type deleteModBody struct {
		ModId string `json:"modId"`
	}

	var parsedBody deleteModBody

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jlog.Errorf("Unable to read \"/admin/addMod\" body.\n")
		jlog.Errorf("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message": "Unable to read body."}`)
		return 
	}

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		jlog.Warnf("Unable to unmarshal JSON data.\n")
		jlog.Println("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message":"Unable to parse body."}`)
		return
	}


	err = app.ModHolder.DeleteMod(parsedBody.ModId)

	if err != nil {
		fmt.Fprintf(w, `{"status": "fail", "message":"` + err.Error() + `"}`)
		return
	}

	fmt.Fprintf(w, `{"status": "good","message":"Module deleted."}`)
}

func getModConfig(app *MainApp, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	

	type getModConfigBody struct {
		ModId string `json:"modId"`
	}

	var parsedBody getModConfigBody

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jlog.Errorf("Unable to read \"/admin/addMod\" body.\n")
		jlog.Errorf("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message": "Unable to read body."}`)
		return 
	}

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		jlog.Warnf("Unable to unmarshal JSON data.\n")
		jlog.Println("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message":"Unable to parse body."}`)
		return
	}

	if mod, ok := app.ModHolder.Mods[parsedBody.ModId]; ok {
		modConfigStr, err := mod.ConfigStr()
		if err != nil {
			jlog.Errorf("%v\n", err)
			fmt.Fprintf(w, `{"status": "fail", "message":"Unable get config string."}`)
			return 
		}

		fmt.Fprintf(w, `{"status":"good","modConfig":` + string(modConfigStr) + `}`)
	} else {
		jlog.Errorf("Unable to retrieve module \"%s\".", parsedBody.ModId)
		fmt.Fprintf(w, `{"status": "fail", "message":"Unable retrieve module"}`)
	}
}

func updateMod(app *MainApp, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	

	type updateModBody struct {
		ModId string `json:"modId"`
		ConfigStr string `json:'configStr'`
	}	

	var parsedBody updateModBody

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jlog.Errorf("Unable to read \"/admin/addMod\" body.\n")
		jlog.Errorf("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message": "Unable to read body."}`)
		return 
	}

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		jlog.Warnf("Unable to unmarshal JSON data.\n")
		jlog.Println("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message":"Unable to parse body."}`)
		return
	}

	err = app.ModHolder.UpdateMod(parsedBody.ModId, parsedBody.ConfigStr)
	if err != nil {
		jlog.Errorf("Unable to update \"%s\" config.\n", parsedBody.ModId)
		jlog.Errorf("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message":"` + err.Error() + `"}`)
		return
	}

	err = app.SyncConfig(parsedBody.ModId)
	if err != nil {
		jlog.Errorf("Unable to sync config to jablkoconfig.json.\n")	
		jlog.Errorf("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message":"` + err.Error() + `"}`)
		return 
	}

	fmt.Fprintf(w, `{"status":"good","message":"Updated config."}`)
}

func addUser(app *MainApp, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	

	type addUserBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		FirstName string `json:"firstName"`
	}

	var parsedBody addUserBody

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jlog.Errorf("Unable to read \"/admin/addMod\" body.\n")
		jlog.Errorf("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message": "Unable to read body."}`)
		return 
	}

	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		jlog.Warnf("Unable to unmarshal JSON data.\n")
		jlog.Println("%v\n", err)
		fmt.Fprintf(w, `{"status": "fail", "message":"Unable to parse body."}`)
		return
	}

	jlog.Println(parsedBody)

	err = app.Db.AddUser(parsedBody.Username, parsedBody.Password, parsedBody.FirstName, 0)
	if err != nil {
		jlog.Warnf("Unable to add to SQLite database\n")
		jlog.Println("%v\n", err)
		fmt.Fprintf(w, `{"status":"fail","message":"` + err.Error() + `"}`)	
		return
	}

	fmt.Fprintf(w, `{"status":"good","message":"Added user."}`)	
}
