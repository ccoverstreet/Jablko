// jablkomodules.go: Jablko Module Manager

package jablkomods

import (
	"fmt"
	"os"
	"plugin"

	"jablko/types"

	"github.com/buger/jsonparser"
)

var ModMap = make(map[string]types.JablkoMod)

func Initialize(jablkoModConfig []byte) {
	fmt.Printf("%s\n", jablkoModConfig)	

	// Iterate through the JSON object to initialize all instances
	jsonparser.ObjectEach(jablkoModConfig, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		// Checks if Jablko Module Package was initialized with 
		// a compiled plugin. If jablkomod.so not found, Jablko
		// will attempt to build the plugin.

		fmt.Printf("%s\n%s\n", key, value)
		
		// DEV WARNING: ONLY WORKS FOR LOCAL MODULES WITH ABSOLUTE PATH
		pluginDir, err := jsonparser.GetString(value, "source")
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		pluginFile := pluginDir + "/jablkomod.so"

		// Check if the plugin has already been built
		if _, err := os.Stat(pluginFile); os.IsNotExist(err) {
			fmt.Printf("Plugin file not found.\n")
		}

		// Load plugin
		plug, err := plugin.Open(pluginFile)	
		if err != nil {
			return err
		}

		// Look for Initialize function symbol in plugin
		initSym, err := plug.Lookup("Initialize")
		if err != nil {
			return err
		}

		// Check if function signature matches
		initFunc, ok := initSym.(func(string, []byte)(types.JablkoMod, error))
		if !ok {
			return nil
		}

		modInstance, err := initFunc(string(key), value)
		if err != nil {
			return err
		}

		ModMap[string(key)] = modInstance

		return nil
	}) 
}
