// Jablko Module Communication
// Cale Overstreet
// February 2, 2021
// Responsible for registering modules and storing module name, auth string, and ip address

package modcommunication

import (
	"sync"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"

	"github.com/ccoverstreet/Jablko/src/jlog"
)

type registryEntry struct {
	Address string `json:"address"`
	AuthStr string `json:"authStr"`
}

type ModRegistry struct {
	sync.Mutex	
	mods map[string]registryEntry
}

func (reg *ModRegistry) InitializeRegistry() error {
	reg.Lock()
	defer reg.Unlock()

	reg.mods = make(map[string]registryEntry)

	if _, err := os.Stat("./data/modregistry.json"); os.IsNotExist(err) {
		jlog.Warnf("Module registry not found in \"data\" directory.\n")
		jlog.Warnf("Initializing empty registry.\n")

		return nil
	} else if err != nil {
		return err
	}

	regData, err := ioutil.ReadFile("./data/modregistry.json")
	if err != nil {
		jlog.Errorf("Unable to retrieve registry data.\n")
		panic(err)
	}

	err = json.Unmarshal(regData, &reg.mods)
	if err != nil {
		panic(err)
	}

	return nil
}

func (reg *ModRegistry) SaveConfig() error {
	saveData, err := json.MarshalIndent(reg.mods, "", "\t")	
	if err != nil {
		panic(err)
	}

	jlog.Println(string(saveData))

	err = ioutil.WriteFile("./data/modregistry.json", saveData, 0022)

	return err
}

func (reg *ModRegistry) AddDevice(name string, address string, authStr string) error {
	jlog.Println(name, address, authStr)

	if _, ok := reg.mods[name]; ok {
		return fmt.Errorf("Module with the name \"%s\" already exists.", name)
	}

	reg.mods[name] = registryEntry{address, authStr}

	reg.SaveConfig()
	return nil	
}
