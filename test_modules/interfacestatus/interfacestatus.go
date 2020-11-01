package main

import (
	"github.com/ccoverstreet/Jablko/types"
	"net/http"
	"log"
	"github.com/buger/jsonparser"
	"strings"
	"strconv"
)

type intStatus struct {
	id string
	title string
	updateInterval int
	jablko types.JablkoInterface
}

func Initialize(instanceId string, configData []byte, jablko types.JablkoInterface) (types.JablkoMod, error) {
	instance := new(intStatus) 

	instance.id = instanceId

	updateInt, err := jsonparser.GetInt(configData, "updateInterval")
	if err != nil {
		return nil, err
	}

	instance.updateInterval = int(updateInt)

	configTitle, err := jsonparser.GetString(configData, "title")
	if err != nil {
		return nil, err
	}

	instance.title = configTitle

	instance.jablko = jablko

	return types.StructToMod(instance), nil
}

func (instance *intStatus) Card(*http.Request) string {
	r := strings.NewReplacer("$UPDATE_INTERVAL", strconv.Itoa(instance.updateInterval),
	"$MODULE_ID", instance.id,
	"$MODULE_TITLE", instance.title)

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
		instance.jablko.Tester()
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

