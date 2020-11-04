package main

import (
	"github.com/ccoverstreet/Jablko/types"
	"net/http"
	"log"
	"github.com/buger/jsonparser"
	"strings"
	"strconv"
	"encoding/json"
)

type intStatus struct {
	id string
	Title string
	UpdateInterval int
	jablko types.JablkoInterface
}

func Initialize(instanceId string, configData []byte, jablko types.JablkoInterface) (types.JablkoMod, error) {
	instance := new(intStatus) 

	instance.id = instanceId

	updateInt, err := jsonparser.GetInt(configData, "UpdateInterval")
	if err != nil {
		return nil, err
	}

	instance.UpdateInterval = int(updateInt)

	configTitle, err := jsonparser.GetString(configData, "Title")
	if err != nil {
		return nil, err
	}

	instance.Title = configTitle

	instance.jablko = jablko

	return types.StructToMod(instance), nil
}

func (instance *intStatus) ConfigStr() ([]byte, error) {
	res, err := json.Marshal(instance)	
	if err != nil {
		return nil, nil	
	}

	log.Printf("%s\n", res)

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
		instance.jablko.Tester()

		go func() {
			x := 0

			for i := 0; i < 10000000; i++ {
				x = i + i
			}

			log.Println(x)
		}()

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

