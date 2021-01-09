// jablkomods.go: Jablko Module Manager

package jablkomods

import (
	"fmt"
	"os"
    "os/exec"
    "syscall"
	"plugin"
	"encoding/json"
	"strings"
	"sync"

	"github.com/ccoverstreet/Jablko/types"
	"github.com/ccoverstreet/Jablko/src/jlog"

	"github.com/buger/jsonparser"
)

type JablkoModuleHolder struct {
	sync.Mutex
	Mods map[string]types.JablkoMod
	Config map[string]string
	Order []string
	mainInterface types.JablkoInterface
}

func Initialize(jablkoModConfig []byte, moduleOrder []byte, jablko types.JablkoInterface) (*JablkoModuleHolder, error) {
	x := new(JablkoModuleHolder)
	x.Mods = make(map[string]types.JablkoMod)
	x.Config = make(map[string]string)
	x.mainInterface = jablko

	flagBuildAll := jablko.GetFlagValue("--build-all")

	// Build cache keeps track of what plugins have been built
	// Only important if --build-all flag is passed
	buildCache := make(map[string]bool)

	// Get the module order
	err := json.Unmarshal(moduleOrder, &x.Order)
	if err != nil {
		jlog.Errorf("Unable to unmarshal module order: %v\n", err)
		panic(err)
	}

	flagRestart := false

	initErr := jsonparser.ObjectEach(jablkoModConfig, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
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

		// Get init func
		initFunc, err := GetPluginInitFunc(pluginFile)
		if err != nil {
			jlog.Errorf("Error getting jablkomod init function.\n")
			jlog.Errorf("%v\n", err)

			// Check if error is fixable
			if strings.Contains(err.Error(), "plugin was built with a different version") || os.IsNotExist(err) {
				jlog.Warnf("Attempting to rebuild \"%s\"...\n", pluginFile)

				// Attempt module rebuild
				err = BuildJablkoMod(installDir)
				if err != nil {
					jlog.Errorf("Unable to rebuild \"%s\".\n")
					jlog.Warnf("Plugin file \"%s\" not found. Jablkomod will not be enabled\n", pluginFile)
					jlog.Warnf("%v\n", err)
					return nil
				}

				jlog.Warnf("Rebuilt plugin \"%s\". Jablko will restart after initialization.\n", pluginFile)
				flagRestart = true
				return nil
			}

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

	if flagRestart {
		binary, lookErr := exec.LookPath("go")
		if lookErr != nil {
			jlog.Errorf("Unable to automatically restart Jablko.\n")
			jlog.Warnf("%v\n", lookErr)
		}

		args := []string{"go", "run", "jablko.go"}

		env := os.Environ()

		execErr := syscall.Exec(binary, args, env)
		if execErr != nil {
			jlog.Errorf("Unable to automatically restart Jablko.\n")
			panic(execErr)
		}
	}

	return x, initErr
}

func (instance *JablkoModuleHolder) InstallMod(modPath string) error {
	instance.Lock()
	defer instance.Unlock()

	splitPath := strings.Split(modPath, "/")
	splitRepo := strings.Split(splitPath[len(splitPath) - 1], "-")

	// Check if source is already present
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		// Download source and build
		err = DownloadJablkoMod(modPath)

		if err != nil {
			jlog.Errorf("Unable to install Jablko Mod.\n")
			return err
		}
	}

	modId := splitRepo[0] + CreateUUID()

	// Get pointer to plugin and calll initialize function
	initFunc, err := GetPluginInitFunc(modPath + "/jablkomod.so")
	if err != nil {
		return err
	}

	newMod, err := initFunc(modId, nil, instance.mainInterface)
	if err != nil {
		return err
	}


	instance.Mods[modId] = newMod
	instance.Order = append(instance.Order, modId)
	instance.UpdateMod(modId, `{"Source": "` + modPath + `"}`)
	instance.mainInterface.SyncConfig(modId)

	return nil
}

func (instance *JablkoModuleHolder) DeleteMod(modId string) error {
	instance.Lock()
	defer instance.Unlock()

	// Check if mod exists in ModHolder.Mods and delete
	if _, ok := instance.Mods[modId]; ok {
		jlog.Printf("Deleting \"%s\" from Mods.\n", modId)
		delete(instance.Mods, modId)		
	}

	// Check if mod exists in ModHolder.Config and delete
	if _, ok := instance.Config[modId]; ok {
		jlog.Printf("Deleting \"%s\" from Config.\n", modId)
		delete(instance.Config, modId)
	}

	// Remove mod from modOrder
	modIndex := -1
	for i := 0; i < len(instance.Order); i++ {
		if instance.Order[i] == modId {
			modIndex = i
			break
		}
	}

	// Delete from modOrder if index > 0
	if modIndex >= 0 {
		instance.Order = append(instance.Order[:modIndex], instance.Order[modIndex + 1:]...)
	} else {
		return fmt.Errorf("Unable to delete mod \"%s\" from module order.", modId)
	}


	instance.mainInterface.SyncConfig("")
	
	return nil
}

func (instance *JablkoModuleHolder) UpdateMod(modId string, newConfig string) error {
	if mod, ok := instance.Mods[modId]; ok {
		return mod.UpdateConfig([]byte(newConfig))
	} else {
		return fmt.Errorf("Unable to update mod config.")
	}
}

func GetPluginInitFunc(pluginFile string) (func(string, []byte, types.JablkoInterface)(types.JablkoMod, error), error) {
	// Check if the plugin has already been built
	if _, err := os.Stat(pluginFile); os.IsNotExist(err) {
		return nil, err
	}

	// Load plugin
	plug, err := plugin.Open(pluginFile)	
	if err != nil {
		return nil,err
	}

	initSym, err := plug.Lookup("Initialize")
	if err != nil {
		return nil, err
	}

	initFunc, ok := initSym.(func(string, []byte, types.JablkoInterface)(types.JablkoMod, error))
	if !ok {
		return nil, fmt.Errorf("Function signature does not match.")
	}

	return initFunc, nil
}
