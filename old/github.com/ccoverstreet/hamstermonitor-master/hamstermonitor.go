package main

import (
	"fmt"
	"strings"
	"time"
	"strconv"
	"sync"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"

	"github.com/ccoverstreet/Jablko/types"
	"github.com/ccoverstreet/Jablko/src/jlog"
)

// ---------- Module Globals ----------
// Interface reference to MainApp
var jablko types.JablkoInterface
var templateCaching bool
var cachedTemplate string

const defaultSourcePath = "github.com/ccoverstreet/hamstermonitor-master" // Defaults to master branch

// ---------- END Module Globals ----------

const activeLength = 240
const storageBin = 2880

type hamsterMonitor struct {
	sync.Mutex
	id string
	Title string
	Source string
	HamsterName string
	active int
	lastActive int64 
	storage [storageBin]byte
	storageTime [storageBin]int64
	storageCounter int
}


func Initialize(instanceId string, configData []byte, jablkoRef types.JablkoInterface) (types.JablkoMod, error) {
	instance := new(hamsterMonitor)
	instance.id = instanceId

	if configData == nil {
		instance.Title = "Hamster Monitor"
		instance.Source = defaultSourcePath
		instance.HamsterName = "Pip"

		return instance, nil
	}

	err := json.Unmarshal(configData, &instance)
	if err != nil {
		return nil, err		
	}


	jablko = jablkoRef
	jablko.SendMessage("Hamster monitor started")

	templateCaching = !jablko.GetFlagValue("--debug-mode")

	if templateCaching && cachedTemplate == "" {
		// Load in the template if not loaded
		templateBytes, err := ioutil.ReadFile(instance.Source + "/hamstermonitor.html")
		if err != nil {
			jlog.Errorf("Unable to read hamstermonitor.html\n")
		}

		cachedTemplate = string(templateBytes)
	}

	return types.StructToMod(instance), nil
}

func (instance *hamsterMonitor) ConfigStr() ([]byte, error) {
	res, err := json.Marshal(instance);

	if err != nil {
		return nil, err
	}

	jlog.Println(instance)

	return res, nil
}

func (instance *hamsterMonitor) UpdateConfig(newConfig []byte) error {
	err := json.Unmarshal(newConfig, instance)
	if err != nil {
		return err
	}

	return nil
}

func (instance *hamsterMonitor) Card(*http.Request) string {
	instance.Lock()
	defer instance.Unlock()

	r := strings.NewReplacer("$MODULE_ID", instance.id,
		"$MODULE_TITLE", instance.Title,
		"$UPDATE_INTERVAL", strconv.Itoa(10), 
		"$HAMSTER_NAME", instance.HamsterName)

	if templateCaching {
		return r.Replace(cachedTemplate)
	}

	templateBytes, err := ioutil.ReadFile(instance.Source + "/hamstermonitor.html")
	if err != nil {
		jlog.Errorf("Unable to read hamstermonitor.html\n")
	}

	htmlTemplate := string(templateBytes)
	return r.Replace(htmlTemplate)
}

type monitorData struct {
	Active int `json:"active"`
}

func (instance *hamsterMonitor) WebHandler(w http.ResponseWriter, r *http.Request) {		
	pathParams := mux.Vars(r)

	if pathParams["func"] == "dump" {
		instance.dataDump(w, r)	
		return
	} else if pathParams["func"] == "getStatus" {
		instance.sendStatus(w, r)
	}
}

func (instance *hamsterMonitor) dataDump(w http.ResponseWriter, r *http.Request) {
	instance.Lock()
	defer instance.Unlock()
	var newData monitorData

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jlog.Errorf("%v\n", err)
	}

	err = json.Unmarshal(body, &newData)
	if err != nil {
		jlog.Errorf("%v\n", err)
	}

	instance.active = newData.Active
	if instance.active == 1 {
		instance.lastActive = time.Now().Unix()
	}

	if instance.storageCounter >= storageBin {
		instance.storageCounter = 0
	}

	instance.storage[instance.storageCounter] = byte(instance.active) + '0'
	instance.storageTime[instance.storageCounter] = time.Now().Unix()
	instance.storageCounter = instance.storageCounter + 1

	jlog.Println(instance.storage)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status": "fail", "message": "Unable to find an appropriate action."}`)
}

func (instance *hamsterMonitor) sendStatus(w http.ResponseWriter, r *http.Request) {
	instance.Lock()
	defer instance.Unlock()

	curActive := 0
	
	if (time.Now().Unix() - instance.lastActive < activeLength) {
		curActive = 1
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status": "fail", "active": ` + strconv.Itoa(curActive) + `}`)
}
