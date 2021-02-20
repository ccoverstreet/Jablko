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

func (instance *intStatus) SourcePath() string {
	return instance.Source
}

func (instance *intStatus) UpdateConfig(newConfig []byte) error {
	err := json.Unmarshal(newConfig, instance)
	if err != nil {
		return err
	}

	return nil
}

func (instance *intStatus) ModuleCardConfig() string {
	type configPayload struct {
		Id string `json:"id"`
		Title string `json:"title"`
		UpdateInterval int `json:"updateInterval"`
	}

	structData := configPayload{instance.id, instance.Title, instance.UpdateInterval}

	data, err := json.Marshal(structData)
	if err != nil {
		jlog.Warnf("builtin/interfacestatus: Unable to marshal module card config to string\n")
	}

	return string(data) 
}

func (instance *intStatus) WebHandler(w http.ResponseWriter, r *http.Request) {
	// Use mux.Vars(r) to route incoming requests
	pathParams := mux.Vars(r)

	var err error = nil

	switch {
	case pathParams["func"] == "getStatus":
		err = getStatus(w, r)	
	case pathParams["func"] == "speedTest":
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"good","message":"Speed test succesful"}`)
	default:
		err = fmt.Errorf("No corresponding function found.")
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status": "fail", "message": "` + err.Error() + `"}`)
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
