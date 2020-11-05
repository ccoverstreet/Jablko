package main

import (
	"github.com/ccoverstreet/Jablko/types"
	"net/http"
	"log"
	"strings"
	"strconv"
	"encoding/json"
)

var jablko types.JablkoInterface

type intStatus struct {
	id string
	Title string
	Source string
	UpdateInterval int
}

func Initialize(instanceId string, configData []byte, jablkoRef types.JablkoInterface) (types.JablkoMod, error) {
	instance := new(intStatus) 

	err := json.Unmarshal(configData, &instance)	
	if err != nil {
		return nil, err		
	}

	log.Println(instance)
	instance.id = instanceId

	jablko = jablkoRef

	return types.StructToMod(instance), nil
}

func (instance *intStatus) ConfigStr() ([]byte, error) {
	instance.UpdateInterval = 5
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

	return r.Replace(htmlTemplate)
}

func (instance *intStatus) WebHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	splitPath := strings.Split(r.URL.Path, "/")
	log.Println(splitPath)
	if len(splitPath) != 4 {
		// Incorrect path received
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid path received."}`))

		return
	}

	switch {
	case splitPath[3] == "fart":
		log.Println("Fart was called by client")
		jablko.Tester()

		jablko.SyncConfig(instance.id)
	case splitPath[3] == "getStatus":
		log.Println("Get status called")
	default:
		log.Println("No call found.")	
	}
}

const htmlTemplate = `
<script>
	const $MODULE_ID = {
		"warn": function() {
			fetch("/jablkomods/$MODULE_ID/fart", {
				method: "POST"
			})
				.then(async (data) => {
					console.log(await data.json())
				})
		}
	}
</script>
<div class="module_card">
	<h1>$MODULE_TITLE</h1>
	<p>Module ID: $MODULE_ID</p>
	<p>Update Interval: $UPDATE_INTERVAL s</p>
</div>
`

