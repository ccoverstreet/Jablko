// jablkomods.go: Jablko Module Manager

package jablkomods

import (
	"fmt"
	"os"
	"plugin"
	"encoding/json"
	"strings"

	"github.com/ccoverstreet/Jablko/types"
	"github.com/ccoverstreet/Jablko/src/jlog"

	"github.com/buger/jsonparser"
)

type JablkoModuleHolder struct {
	Mods map[string]types.JablkoMod
	Config map[string]string
	Order []string
}

func Initialize(jablkoModConfig []byte, moduleOrder []byte, jablko types.JablkoInterface) (*JablkoModuleHolder, error) {
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
		jlog.Errorf("ERROR: Unable to unmarshal module order.")
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

		// Check if package needs to be downloaded
		if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
			jlog.Printf("Unable to find jablkomod directory %s.\n", pluginDir)	
			if strings.HasPrefix(pluginDir, "github.com") {
				jlog.Printf("Downloading plugin source\n")
				downloadErr := DownloadJablkoMod(pluginDir)

				if downloadErr != nil {
					jlog.Errorf("Unable to download plugin \"%s\"\n", pluginDir)
					jlog.Errorf("%v\n", downloadErr)
					jlog.Warnf("Module \"%s\" will not be enabled\n", pluginDir)
					return nil // THIS IS KIND OF STRANGE
					// Currently returns nil so all modules are parsed
				}
			} else {
				jlog.Warnf("WARNING: Jablkomod %s not found and will not be downloaded. Jablkomod will not be enabled.\n", pluginDir)
			}
		}

		installDir := pluginDir

		// Build plugin if flagBuildAll is true
		if flagBuildAll {
			jlog.Printf("Checking if jablkomod \"%s\" needs to be built.\n", installDir)
			if _, ok := buildCache[installDir]; !ok {
				jlog.Printf("Building jablkomod \"%s\".\n", installDir)
				err = BuildJablkoMod(installDir)
				jlog.Println(installDir)
				if err != nil {
					jlog.Errorf("ERROR: Unable to build jablkomod located in \"%s\".\n", installDir)
					jlog.Warnf("WARNING: Jablko %s will not be activated.\n", installDir)
				}
				buildCache[installDir] = true
			} else {
				jlog.Printf("Jablkomod already built.\n")
			}
		}

		pluginFile := installDir + "/jablkomod.so"

		// Check if the plugin has already been built
		if _, err := os.Stat(pluginFile); os.IsNotExist(err) {
			jlog.Warnf("Plugin file \"%s\" not found. Jablkomod will not be enabled\n", pluginFile)
			jlog.Warnf("%v\n", err)
			return nil
		}

		// Load plugin
		plug, err := plugin.Open(pluginFile)	
		if err != nil {
			jlog.Warnf("Error loading jablkomod \"%s\".\n", pluginFile)
			jlog.Warnf("%v\n", err)

			if strings.Contains(err.Error(), "plugin was built with a different version") {
				jlog.Warnf("Attempting to rebuild \"%s\"\n...", pluginFile)

				// Attempt module rebuild
				err = BuildJablkoMod(installDir)
				if err != nil {
					jlog.Errorf("Unable to rebuild \"%s\".\n")
					jlog.Warnf("Plugin file \"%s\" not found. Jablkomod will not be enabled\n", pluginFile)
					jlog.Warnf("%v\n", err)
					return nil
				}
			}
		}

		// Look for Initialize function symbol in plugin
		initSym, err := plug.Lookup("Initialize")
		if err != nil {
			jlog.Warnf("Initialize function signature not found in \"%s\"\n", pluginFile)
			jlog.Warnf("%v\n", err)
			return nil
		}

		// Check if function signature matches
		initFunc, ok := initSym.(func(string, []byte, types.JablkoInterface)(types.JablkoMod, error))
		if !ok {
			jlog.Warnf("Initialize function signature doesn't match \"%s\"\n", pluginFile)
			return nil
		}

		modInstance, err := initFunc(string(key), value, jablko)
		if err != nil {
			jlog.Warnf("Initialize function failed in \"%s\"\n", pluginFile)
			return nil
		}

		x.Mods[string(key)] = modInstance

		return nil
	})
}

func (instance *JablkoModuleHolder) InstallMod(modPath string) error {
	jlog.Println(GithubSourceToURL(modPath))					
		
	return nil
}
