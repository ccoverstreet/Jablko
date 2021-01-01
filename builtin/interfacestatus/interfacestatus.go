package main

import (
	"github.com/ccoverstreet/Jablko/types"
	"github.com/ccoverstreet/Jablko/src/jlog"

	"net/http"
	"fmt"
	"time"
	"runtime"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
)

// ---------- Module Globals ----------
var jablko types.JablkoInterface
var templateCaching bool
var cachedTemplate string

var serverStartTime int
// ---------- END Module Globals ----------

type intStatus struct {
	id string
	Title string
	Source string
	UpdateInterval int
}

func init() {
	// Initialiaze globals
	serverStartTime = int(time.Now().Unix())
}

func Initialize(instanceId string, configData []byte, jablkoRef types.JablkoInterface) (types.JablkoMod, error) {
	instance := new(intStatus) 
	instance.id = instanceId

	// Return default config if no configData is supplied
	if configData == nil {
		instance.Title = "Interface Status"
		instance.Source = "./builtin/interfacestatus"
		instance.UpdateInterval = 25

		return instance, nil
	}

	// Initialize instance with configData
	err := json.Unmarshal(configData, &instance)	
	if err != nil {
		return nil, err		
	}

	jablko = jablkoRef
	
	templateCaching = !jablko.GetFlagValue("--debug-mode")

	if templateCaching {
		loadedTemplateBytes, err := ioutil.ReadFile(instance.Source + "/interfacestatus.html")
		if err != nil {
			jlog.Errorf("ERROR: Unable to read interfacestatus.html template file\n")
		}

		cachedTemplate = string(loadedTemplateBytes)
	}
	
	return types.StructToMod(instance), nil
}

func (instance *intStatus) ConfigStr() ([]byte, error) {
	res, err := json.Marshal(instance)	
	if err != nil {
		return nil, err
	}

	return res, nil
}


func (instance *intStatus) Card(*http.Request) string {
	r := strings.NewReplacer("$UPDATE_INTERVAL", strconv.Itoa(instance.UpdateInterval),
	"$MODULE_ID", instance.id,
	"$MODULE_TITLE", instance.Title)

	if templateCaching {
		return r.Replace(cachedTemplate)
	}

	loadedTemplateBytes, err := ioutil.ReadFile(instance.Source + "/interfacestatus.html")
	if err != nil {
		jlog.Errorf("ERROR: Unable to read interfacestatus.html template file\n")
	}

	htmlTemplate := string(loadedTemplateBytes)

	return r.Replace(htmlTemplate)
}

func (instance *intStatus) WebHandler(w http.ResponseWriter, r *http.Request) {
	// Use mux.Vars(r) to route incoming requests
	pathParams := mux.Vars(r)

	var err error = nil

	switch {
	case pathParams["func"] == "banana":
		jlog.Println("ASDASDASDSA")
	case pathParams["func"] == "getStatus":
		err = getStatus(w, r)	
	case pathParams["func"] == "speedTest":
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"good","message":"Speed test succesful"}`)
	default:
		jlog.Println("Nothing Found")
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status": "fail", "message": "Unable to find an appropriate action."}`)
	}
}

func getStatus(w http.ResponseWriter, r *http.Request) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	resTemplate := `{"status": "$STATUS", "message": "$MESSAGE", "uptime": $UPTIME, "curAlloc": $CUR_ALLOC, "sysAlloc": $SYS_ALLOC}`

	replacer := strings.NewReplacer("$STATUS", "good",
		"$MESSAGE", "Status normal.",
		"$UPTIME", strconv.Itoa(int(time.Now().Unix()) - serverStartTime),
		"$CUR_ALLOC", strconv.Itoa(int(m.Alloc)),
		"$SYS_ALLOC", strconv.Itoa(int(m.Sys)))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, replacer.Replace(resTemplate))

	return nil
}
