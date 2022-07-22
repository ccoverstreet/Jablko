// HTTP Handlers StatusApp
// Cale Overstreet
//

package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func WrapHandler(f func(w http.ResponseWriter, r *http.Request, inst *StatusApp), inst *StatusApp) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, inst)
	}
}

func ReadJSONBody(body io.ReadCloser, dest interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	log.Println(string(b))

	return json.Unmarshal(b, dest)
}

func InstanceDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[{}]")
}

func WebComponentHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./webcomponent.js")
	if err != nil {
		return
	}

	fmt.Fprintf(w, "%s", b)
}

func RemoveDeviceHandler(w http.ResponseWriter, r *http.Request, inst *StatusApp) {
	type format struct {
		IPAddress string `json:"ipAddress"`
	}

	data := format{}

	ReadJSONBody(r.Body, &data)
	log.Println(data)

	err := inst.RemoveDevice(data.IPAddress)
	if err != nil {
		log.Printf("ERROR: Unable to remove connection - %v", err)
		fmt.Fprintf(w, "%v", err)
		return
	}

	inst.UpdateSummary()
	inst.PushConnections()

	fmt.Fprintf(w, "Removed device from poll list")
}

func AddDeviceHandler(w http.ResponseWriter, r *http.Request, inst *StatusApp) {
	type format struct {
		IPAddress string `json:"ipAddress"`
		Name      string `json:"name"`
	}

	data := format{}

	ReadJSONBody(r.Body, &data)

	log.Println("ASDA", data)

	err := inst.AddDevice(data.IPAddress, data.Name)
	if err != nil {
		log.Printf("ERROR: Unable to add connection - %v", err)
		fmt.Fprintf(w, "%v", err)
		return
	}

	inst.UpdateSummary()
	inst.PushConnections()

	fmt.Fprintf(w, "Added device to poll list")
}
