// jablkomods.go: Jablko Module Manager

package jablkomods

import (
	"fmt"
	"os"
	"plugin"
	"encoding/json"
	"os/exec"
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

var ModMap = make(map[string]types.JablkoMod)

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
				downloadErr := downloadJablkoMod(pluginDir)

				if downloadErr != nil {
					jlog.Errorf("Unable to download plugin \"%s\"\n", pluginDir)
					return nil // THIS IS KIND OF STRANGE
					// Currently returns nil so all modules are parsed
					// Is it worth allowing Jablko to start if not
					// all mods are able to load? Either way, this should 
					// fail loudly.
				}
			} else {
				jlog.Warnf("WARNING: Jablkomod %s not found and will not be downloaded. Jablkomod will not be enabled.\n", pluginDir)
			}
		}

		// Build plugin if flagBuildAll is true
		if flagBuildAll {
			jlog.Printf("Checking if jablkomod \"%s\" needs to be built.\n", pluginDir)
			if _, ok := buildCache[pluginDir]; !ok {
				jlog.Printf("Building jablkomod \"%s\".\n", pluginDir)
				err = buildJablkoMod(pluginDir)
				if err != nil {
					jlog.Errorf("ERROR: Unable to build jablkomod located in \"%s\".\n", pluginDir)
					jlog.Warnf("WARNING: Jablko %s will not be activated.\n", pluginDir)
				}
				buildCache[pluginDir] = true
			} else {
				jlog.Printf("Jablkomod already built.\n")
			}
		}

		pluginFile := pluginDir + "/jablkomod.so"

		// Check if the plugin has already been built
		if _, err := os.Stat(pluginFile); os.IsNotExist(err) {
			jlog.Warnf("Plugin file \"%s\"not found. Jablkomod will not be enabled\n", pluginFile)
			return err
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

func buildJablkoMod(buildDir string) error {
	buildCMD := exec.Command("go", "build", "-buildmode=plugin", "-o", "jablkomod.so", ".")		
	buildCMD.Dir = buildDir
	jlog.Println(buildCMD)
	_, err := buildCMD.Output()

	return err
}

func downloadJablkoMod(repoPath string) error {
	return fmt.Errorf("Unable to download jablkomod repo \"%s\"\n", repoPath)
}
