package main

import (
	"github.com/ccoverstreet/Jablko/types"
	"github.com/ccoverstreet/Jablko/src/jlog"

	"encoding/json"
	"fmt"
	"sync"

	"net/http"
	"github.com/gorilla/mux"
)

// Globals
var jablko types.JablkoInterface
const defaultSourcePath = "./builtin/interfacestatus"

type notepad struct {
	sync.Mutex
	id string
	Title string
	Source string
}

func Initialize(instanceId string, configData []byte, jablkoRef types.JablkoInterface) (types.JablkoMod, error) {
	instance := new(notepad)
	instance.id = instanceId

	if configData == nil {
		// Default
		instance.Title = "Notepad"
		instance.Source = defaultSourcePath

		return instance, nil
	}

	// Initialize with config data
	err := json.Unmarshal(configData, &instance)
	if err != nil {
		return nil, err
	}

	jablko = jablkoRef

	return types.StructToMod(instance), nil
}

func (instance *notepad) ConfigStr() ([]byte, error) {
	instance.Lock()
	defer instance.Unlock()

	res, err := json.Marshal(instance)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (instance *notepad) SourcePath() string {
	return instance.Source
}

func (instance *notepad) UpdateConfig(newConfig []byte) error {
	instance.Lock()
	defer instance.Unlock()

	err := json.Unmarshal(newConfig, instance)
	if err != nil {
		return err
	}

	return nil
}

func (instance *notepad) ModuleCardConfig() string {
	instance.Lock()
	defer instance.Unlock()

	type configPayload struct {
		Id string `json:"id"`
		Title string `json:"title"`
	}

	structData := configPayload{instance.id, instance.Title}

	data, err := json.Marshal(structData)
	if err != nil {
		jlog.Warnf("builtin/notepad: Unable to marshal module card config to string\n")
		return ""
	}

	return string(data)
}

func (instance *notepad) WebHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)

	var err error = nil

	switch pathParams["func"] {
	case "speak":
		jlog.Warnf("ASDASDASD SPEAK\n")
	case "saveNote":
		err = instance.SaveNote(w, r)
		jlog.Warnf("SAVE NOTE\n")
	default: 
		err = fmt.Errorf("No corresponding function found.")
	}

	if err != nil {
		jlog.Errorf("%v\n", err)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status": "fail", "message": "` + err.Error() + `"}`)
	}
}