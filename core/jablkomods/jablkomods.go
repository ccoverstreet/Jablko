// Jablko Mod Manager
// Cale Overstreet
// Mar. 30, 2021

// Responsible for managing mod state and jablkomod
// processes. 

package jablkomods

import (
	"encoding/json"
	"log"
)

type ModManager struct {
	StateMap map[string]ModData
}

type ModData struct {
	Name string `json:"name"`
	Source string `json:"source"`
	Config interface{} `json:"config"`
}


func NewModManager(config string) (*ModManager, error) {
	x := new(ModManager)
	f := make(map[string]ModData)

	err := json.Unmarshal([]byte(config), &f)
	if err != nil {
		return nil, err
	}

	x.StateMap = f
	log.Println(x)

	b, err := json.Marshal(f["test1"])
	if err != nil {
		return nil, err
	}

	log.Println(string(b))

	return x, nil
}
