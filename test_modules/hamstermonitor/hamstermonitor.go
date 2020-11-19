package main

import (
	"log"
	"fmt"
	"strings"
	"time"
	"strconv"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"

	"github.com/ccoverstreet/Jablko/types"
)

var jablko types.JablkoInterface

const activeLength = 4 * 60

type hamsterMonitor struct {
	id string
	Title string
	Source string
	HamsterName string
	active int
	lastActive int64 
}

func Initialize(instanceId string, configData []byte, jablkoRef types.JablkoInterface) (types.JablkoMod, error) {
	instance := new(hamsterMonitor)

	err := json.Unmarshal(configData, &instance)
	if err != nil {
		return nil, err		
	}

	log.Println(instance)
	instance.id = instanceId

	jablko = jablkoRef

	return types.StructToMod(instance), nil
}

func (instance *hamsterMonitor) ConfigStr() ([]byte, error) {
	res, err := json.Marshal(instance);

	if err != nil {
		return nil, err
	}

	log.Println(instance)

	return res, nil
}

func (instance *hamsterMonitor) Card(*http.Request) string {
	r := strings.NewReplacer("$MODULE_ID", instance.id,
		"$UPDATE_INTERVAL", strconv.Itoa(10), 
		"$HAMSTER_NAME", instance.HamsterName)

	templateBytes, err := ioutil.ReadFile(instance.Source + "/hamstermonitor.html")
	if err != nil {
		log.Println("Unable to read hamstermonitor.html")
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
	var newData monitorData

	log.Println("ASDASDASDAS HAMSTER")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(body, &newData)
	if err != nil {
		log.Println(err)
	}

	instance.active = newData.Active
	instance.lastActive = time.Now().Unix()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status": "fail", "message": "Unable to find an appropriate action."}`)

}

func (instance *hamsterMonitor) sendStatus(w http.ResponseWriter, r *http.Request) {
	curActive := 0
	
	if (instance.active == 1) {
		curActive = 1
	} else if (time.Now().Unix() - instance.lastActive < activeLength) {
		curActive = 1
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status": "fail", "active": ` + strconv.Itoa(curActive) + `}`)
}
