package main

import (
	"github.com/ccoverstreet/Jablko/types"
	"net/http"
	"fmt"
	"github.com/buger/jsonparser"
)

type intStatus struct {
	name string
	updateInterval int
}

func Initialize(instanceName string, configData []byte) (types.JablkoMod, error) {
	fmt.Println("Initializing Interface Status Module\n")
	fmt.Printf("%s\n", configData)

	instance := new(intStatus) 

	instance.name = instanceName

	updateInt, err := jsonparser.GetInt(configData, "updateInterval")
	if err != nil {
		return nil, err
	}

	instance.updateInterval = int(updateInt)

	return types.StructToMod(instance), nil
}

func (instance *intStatus) Card(*http.Request) string {
	return fmt.Sprintf("Hello from Interface Status Module with interval of %d s", instance.updateInterval)
}

func (instance *intStatus) WebHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WEB HANLDER")
}
