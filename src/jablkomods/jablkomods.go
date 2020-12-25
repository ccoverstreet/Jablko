// jablkomods.go: Jablko Module Manager

package jablkomods

import (
	"fmt"
	"os"
	"plugin"
	"encoding/json"
	"os/exec"
	"strings"
	"net/http"
	"io"
	"archive/zip"
	"path/filepath"

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
					jlog.Errorf("%v\n", downloadErr)
					jlog.Warnf("Module \"%s\" will not be enabled\n", pluginDir)
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
			jlog.Warnf("Plugin file \"%s\" not found. Jablkomod will not be enabled\n", pluginFile)
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
	// Download zip from github
	resp, err := http.Get("https://" + repoPath + "/archive/master.zip")
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Bad HTTP status code received\n")
	}

	defer resp.Body.Close()

	// Create file
	err = os.MkdirAll("./tmp/" + repoPath, 0755)
	if err != nil {
		return err
	}

	out, err := os.Create("./tmp/" + repoPath +  "/source.zip")
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	err = os.MkdirAll("./" + repoPath, 0755)
	if err != nil {
		return err
	}

	filenames, err := Unzip("./tmp/" + repoPath +  "/source.zip", "./" + repoPath)
	if err != nil {
		return err
	}

	topLevelDirRaw := strings.Split(filenames[0], "/")
	topLevelDir := topLevelDirRaw[len(topLevelDirRaw) - 1]
	repoPathSplit := strings.Split(repoPath, "/")
	authorDir := repoPathSplit[0] + "/" + repoPathSplit[1]

	// Move directory up one level and rename correctly
	err = os.Rename("./" + repoPath + "/" + topLevelDir, "./" + authorDir + "/" + topLevelDir)
	if err != nil {
		return err
	}

	err = os.RemoveAll("./" + repoPath)
	if err != nil {
		return err
	}

	err = os.Rename("./" + authorDir + "/" + topLevelDir, "./" + repoPath)
	if err != nil {
		return err
	}

	return err
}

func Unzip(src string, dest string) ([]string, error) {
	// From golangcode.com

    var filenames []string

    r, err := zip.OpenReader(src)
    if err != nil {
        return filenames, err
    }
    defer r.Close()

    for _, f := range r.File {

        // Store filename/path for returning and using later on
        fpath := filepath.Join(dest, f.Name)

        // Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return filenames, fmt.Errorf("%s: illegal file path", fpath)
        }

        filenames = append(filenames, fpath)

        if f.FileInfo().IsDir() {
            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        // Make File
        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return filenames, err
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return filenames, err
        }

        rc, err := f.Open()
        if err != nil {
            return filenames, err
        }

        _, err = io.Copy(outFile, rc)

        // Close the file without defer to close before next iteration of loop
        outFile.Close()
        rc.Close()

        if err != nil {
            return filenames, err
        }
    }
    return filenames, nil
}
