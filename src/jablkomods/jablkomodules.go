// jablkomodules.go: Jablko Module Manager

package jablkomods

import (
	"fmt"
	"log"
	"os"
	"plugin"
	"encoding/json"
	"os/exec"

	"github.com/ccoverstreet/Jablko/types"

	"github.com/buger/jsonparser"
)

type JablkoModuleHolder struct {
	Mods map[string]types.JablkoMod
	Config map[string]string
	Order []string
}

var ModMap = make(map[string]types.JablkoMod)

func Initialize2(jablkoModConfig []byte, moduleOrder []byte, jablko types.JablkoInterface) (*JablkoModuleHolder, error) {
	x := new(JablkoModuleHolder)
	x.Mods = make(map[string]types.JablkoMod)
	x.Config = make(map[string]string)

	flagBuildAll := jablko.GetFlagValue("--build-all")

	// Build cache keeps track of what plugins have been built
	// Only important if --build-all flag is passed
	buildCache := make(map[string]bool)

	// Get the module order
	err := json.Unmarshal(moduleOrder, &x.Order)
	if err != nil {
		log.Println("ERROR: Unable to unmarshal module order.")
		panic(err)
	}

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

		// Build plugin if flagBuildAll is true
		if flagBuildAll {
			log.Printf("Checking if jablkomod \"%s\" needs to be built.\n", pluginDir)
			if _, ok := buildCache[pluginDir]; !ok {
				log.Printf("Building jablkomod \"%s\".\n", pluginDir)
				err = buildPlugin(pluginDir)
				if err != nil {
					log.Printf("ERROR: Unable to build jablkomod located in \"%s\".\n", pluginDir)
					panic(err)
				}
				buildCache[pluginDir] = true
			} else {
				log.Printf("Jablkomod already built.\n")
			}
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

func buildPlugin(buildDir string) error {
	buildCMD := exec.Command("go", "build", "-buildmode=plugin", "-o", "jablkomod.so", ".")		
	buildCMD.Dir = buildDir
	log.Println(buildCMD)
	_, err := buildCMD.Output()

	return err
}
