package main

import (
	//"github.com/ccoverstreet/Jablko/types"
	"github.com/ccoverstreet/Jablko/src/jlog"

	"encoding/json"
	"fmt"

	"net/http"
	//"github.com/gorilla/mux"
)

func (instance *notepad) SaveNote(w http.ResponseWriter, r *http.Request) error {
	instance.Lock()
	defer instance.Unlock()

	type saveData struct {
		Text string `json:"text"`
	}

	var parsedBody saveData
	err := json.NewDecoder(r.Body).Decode(&parsedBody)
    if err != nil {
        fmt.Fprintf(w, `{"status":"fail","message":"Unable to parse JSON body."}`)
        return err
    }

	jlog.Println(parsedBody)

	return nil	
}
