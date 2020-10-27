package main

import (
	"github.com/ccoverstreet/Jablko/types"
	"net/http"
	"fmt"
	"github.com/buger/jsonparser"
	"strings"
	"strconv"
)

type intStatus struct {
	id string
	title string
	updateInterval int
}

func Initialize(instanceId string, configData []byte) (types.JablkoMod, error) {
	fmt.Println("Initializing Interface Status Module\n")
	fmt.Printf("%s\n", configData)

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

	return types.StructToMod(instance), nil
}

func (instance *intStatus) Card(*http.Request) string {
	r := strings.NewReplacer("$UPDATE_INTERVAL", strconv.Itoa(instance.updateInterval),
		"$MODULE_ID", instance.id,
		"$MODULE_TITLE", instance.title)

	return r.Replace(htmlTemplate)
}

func (instance *intStatus) WebHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WEB HANLDER")
}

const htmlTemplate = `
<script>
	const $MODULE_ID = {
		"warn": function() {
			alert("WARNING")
		}
	}
</script>
<div class="module_card">
	<h1>$MODULE_TITLE</h1>
	<p>Update Interval: $UPDATE_INTERVAL s</p>
</div>
`

