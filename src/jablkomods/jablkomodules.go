// jablkomodules.go: Jablko Module Manager

package jablkomods

import (
	"fmt"
	"os"
	"plugin"

	"github.com/ccoverstreet/Jablko/types"

	"github.com/buger/jsonparser"
)

type JablkoModuleHolder struct {
	Mods map[string]types.JablkoMod
	Config map[string]string
}

var ModMap = make(map[string]types.JablkoMod)

func Initialize2(jablkoModConfig []byte, jablko types.JablkoInterface) (*JablkoModuleHolder, error) {
	x := new(JablkoModuleHolder)
	x.Mods = make(map[string]types.JablkoMod)
	x.Config = make(map[string]string)

	return x, jsonparser.ObjectEach(jablkoModConfig, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		// Checks if Jablko Module Package was initialized with 
		// a compiled plugin. If jablkomod.so not found, Jablko
		// will attempt to build the plugin.

		x.Config[string(key)] = string(value)


		// DEV WARNING: ONLY WORKS FOR LOCAL MODULES WITH ABSOLUTE PATH
		pluginDir, err := jsonparser.GetString(value, "Source")
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
		initFunc, ok := initSym.(func(string, []byte, types.JablkoInterface)(types.JablkoMod, error))
		if !ok {
			return nil
		}

		modInstance, err := initFunc(string(key), value, jablko)
		if err != nil {
			return err
		}

		x.Mods[string(key)] = modInstance

		return nil
	})
}

func Initialize(jablkoModConfig []byte, jablko types.JablkoInterface) (map[string]string, error) {
	// Iterate through the JSON object to initialize all instances

	x := new(JablkoModuleHolder)
	x.Mods = make(map[string]types.JablkoMod)
	x.Config = make(map[string]string)

	configMap := make(map[string]string)
	
	return configMap, jsonparser.ObjectEach(jablkoModConfig, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		// Checks if Jablko Module Package was initialized with 
		// a compiled plugin. If jablkomod.so not found, Jablko
		// will attempt to build the plugin.

		configMap[string(key)] = string(value)
		x.Config[string(key)] = string(value)


		// DEV WARNING: ONLY WORKS FOR LOCAL MODULES WITH ABSOLUTE PATH
		pluginDir, err := jsonparser.GetString(value, "Source")
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
		initFunc, ok := initSym.(func(string, []byte, types.JablkoInterface)(types.JablkoMod, error))
		if !ok {
			return nil
		}

		modInstance, err := initFunc(string(key), value, jablko)
		if err != nil {
			return err
		}

		x.Mods[string(key)] = modInstance
		ModMap[string(key)] = modInstance

		return nil
	}) 
}
