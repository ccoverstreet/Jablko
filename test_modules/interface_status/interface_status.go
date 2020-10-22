package main

import (
	"jablko/types"
	"net/http"
	"fmt"
	"github.com/buger/jsonparser"
)

type intStatus struct {
	updateInterval int
}

func Initialize(instanceName string, configData []byte) (types.JablkoMod, error) {
	fmt.Println("Initializing Interface Status Module\n")
	fmt.Printf("%s\n", configData)

	instance := new(intStatus) 

	updateInt, err := jsonparser.GetInt(configData, "updateInterval")
	if err != nil {
		return nil, err
	}

	instance.updateInterval = int(updateInt)

	return types.StructToMod(instance), nil
}

func (instance *intStatus) Card(*http.Request) string {
	return "Hello"
}

func (instance *intStatus) WebHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WEB HANLDER")
}
