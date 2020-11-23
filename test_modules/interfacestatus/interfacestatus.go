package main

import (
	"github.com/ccoverstreet/Jablko/types"
	"net/http"
	"fmt"
	"log"
	"time"
	"runtime"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
)

var jablko types.JablkoInterface

var serverStartTime int

type intStatus struct {
	id string
	Title string
	Source string
	UpdateInterval int
}

func init() {
	serverStartTime = int(time.Now().Unix())
	log.Println(serverStartTime)
}

func Initialize(instanceId string, configData []byte, jablkoRef types.JablkoInterface) (types.JablkoMod, error) {
	instance := new(intStatus) 

	err := json.Unmarshal(configData, &instance)	
	if err != nil {
		return nil, err		
	}

	instance.id = instanceId

	jablko = jablkoRef

	return types.StructToMod(instance), nil
}

func (instance *intStatus) ConfigStr() ([]byte, error) {
	res, err := json.Marshal(instance)	
	if err != nil {
		return nil, err
	}

	log.Println(instance)

	return res, nil
}

func (instance *intStatus) Card(*http.Request) string {
	r := strings.NewReplacer("$UPDATE_INTERVAL", strconv.Itoa(instance.UpdateInterval),
	"$MODULE_ID", instance.id,
	"$MODULE_TITLE", instance.Title)

	loadedTemplateBytes, err := ioutil.ReadFile(instance.Source + "/interfacestatus.html")
	if err != nil {
		log.Println("ERROR: Unable to read interfacestatus.html template file")
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
		log.Println("ASDASDASDSA")
	case pathParams["func"] == "getStatus":
		err = getStatus(w, r)	
	default:
		log.Println("Nothing Found")
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
