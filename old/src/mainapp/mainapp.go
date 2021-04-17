// mainapp.go: Main application definitions for Jablko
// Cale Overstreet
// 2020/24/11
// Contains main internal definitions for Jablko

package mainapp

import (
	"strconv"
	"strings"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"fmt"

	"github.com/buger/jsonparser"

	"github.com/ccoverstreet/Jablko/src/jablkomods"
	"github.com/ccoverstreet/Jablko/src/database"
	"github.com/ccoverstreet/Jablko/src/jlog"
	"github.com/ccoverstreet/Jablko/src/modcommunication"
)

type generalConfig struct {
	HttpPort int
	HttpsPort int
	CertPath string
	PrivKeyPath string
}

var jablkoConfig = generalConfig{HttpPort: 8080, HttpsPort: -1}

type MainApp struct {
	sync.Mutex
	Config generalConfig
	ModHolder *jablkomods.JablkoModuleHolder
	Db *database.JablkoDB
	ModRegistry modcommunication.ModRegistry
	flags map[string]bool
}

func CreateMainApp(configFilePath string) (*MainApp, error) {
	// Load config file
	if _, err := os.Stat("./jablkoconfig.json"); os.IsNotExist(err) {
		err = copyDefaultConfig()

		if err != nil {
			jlog.Errorf("Unable to copy default config file. Aborting startup.\n")
			panic (err)
		}
	} else if err != nil {
		jlog.Errorf("Unable to check status of config file. Aborting startup.\n")
		panic(err)
	}

	configData, err := ioutil.ReadFile("./jablkoconfig.json")
	if err != nil {
		jlog.Errorf("%v\n", err)
		panic(err)
	}

	instance := new(MainApp)

	instance.flags = make(map[string]bool)

	// Parse and set flags
	// First set defaults
	for _, val := range os.Args {
		switch val {
		case "--build-all":
			instance.flags[val] = true
		case "--debug-mode":
			instance.flags[val] = true
		}
	}

	httpPort, err := jsonparser.GetInt(configData, "http", "port")
	if err != nil {
		jlog.Errorf("%v\n", err)
		panic("Error getting HTTP port data\n")
	}

	httpsPort, err := jsonparser.GetInt(configData, "https", "port")
	if err != nil {
		jlog.Warnf("HTTPS port Config not set in Config file\n")
	} else {
		jablkoConfig.HttpsPort = int(httpsPort)
	}

	if httpsPort > 0 {
		certPath, err := jsonparser.GetString(configData, "https", "cert")
		if err != nil {
			jlog.Errorf("HTTPS cert path not readable.\n")
			jlog.Errorf("%\v\n", err)
			panic(err)
		}

		privKeyPath, err := jsonparser.GetString(configData, "https", "privKey")
		if err != nil {
			jlog.Errorf("HTTPS private key path not readable.\n")
			jlog.Errorf("%\v\n", err)
			panic(err)
		}

		instance.Config.CertPath = certPath
		instance.Config.PrivKeyPath = privKeyPath
	}

	instance.Config.HttpPort = int(httpPort)
	instance.Config.HttpsPort = int(httpsPort)

	jablkoModulesSlice, _, _, err := jsonparser.Get(configData, "jablkoModules")
	if err != nil {
		panic("Error get Jablko Modules Config\n")
	}

	moduleOrderSlice, _, _, err := jsonparser.Get(configData, "moduleOrder")

	newModHolder, err := jablkomods.Initialize(jablkoModulesSlice, moduleOrderSlice, instance)

	if err != nil {
		jlog.Errorf("JablkoMods ERROR: %s\n", err)
	}

	instance.ModHolder = newModHolder

	instance.Db = database.Initialize()

	// Initialize module registry
	err = instance.ModRegistry.InitializeRegistry()
	if err != nil {
		jlog.Errorf("Unable to initialize module registry.\n")
		jlog.Errorf("%v\n", err)
	}

	err = instance.ModRegistry.AddDevice("RPi" ,"10.0.0.90", "12312hkaksdfaksdjfgasdf")
	if err != nil {
		jlog.Errorf("%v\n", err)
	}


	return instance, nil
}

func (app *MainApp) SendMessage(message string) error {
	jlog.Printf("Message: %s\n", message)	

	return nil
} 

func (app *MainApp) GetFlagValue(flag string) bool {
	if val, ok := app.flags[flag]; ok {
		return val
	} else {
		return false
	}
}

func (app *MainApp) SyncConfig(modId string) error {
	app.Lock()
	defer app.Unlock()

	jlog.Printf("Sync config called for module \"%s\"\n", modId)		

	ConfigTemplate := `{
	"http": {
		"port": $httpPort
	},
	"https": {
		"port": $httpsPort
	},
	"jablkoModules": {
		$moduleString
	},
	"moduleOrder": [
		$moduleOrder
	]
}
`

	if modId == "" {
		// Do nothing, should just dump current state
	} else if _, ok := app.ModHolder.Mods[modId]; !ok {
		jlog.Warnf("Cannot find module %s", modId)
		return fmt.Errorf("Unable to find module.")
	} else {
		newConfByte, err := app.ModHolder.Mods[modId].ConfigStr()
		newConfStr := string(newConfByte)
		if err != nil {
			jlog.Warnf("Unable to get Config string for module %s\n", modId)
			return err
		}

		if app.ModHolder.Config[modId] == newConfStr {
			// If there is no change in config
			return nil
		}

		app.ModHolder.Config[modId] = newConfStr
	}

	// Create JSON to dump to Config file

	// Prepare each module's string
	jablkoModulesStr := ""
	index := 0
	for key, value := range app.ModHolder.Config {
		if index > 0 {
			jablkoModulesStr = jablkoModulesStr + ",\n\t\t\"" + key + "\":" + value
		} else {
			jablkoModulesStr = jablkoModulesStr + "\"" + key + "\":" + value
		}

		index = index + 1
	}

	// Prepare Module Order
	orderStr := ""
	for index, val := range app.ModHolder.Order {
		if index > 0 {
			orderStr = orderStr + ",\n" + "\t\t\"" + val + "\""
		} else {
			orderStr = orderStr + "\"" + val + "\""
		}
	}

	jlog.Println(orderStr)

	r := strings.NewReplacer("$httpPort", strconv.Itoa(jablkoConfig.HttpPort),
	"$httpsPort", strconv.Itoa(jablkoConfig.HttpsPort),
	"$moduleString", jablkoModulesStr,
	"$moduleOrder", orderStr)

	ConfigDumpStr := r.Replace(ConfigTemplate)

	err := ioutil.WriteFile("./jablkoconfig.json", []byte(ConfigDumpStr), 0022)
	return err
}

func copyDefaultConfig() error {
	defaultSrc, err := os.Open("./builtin/defaultconfig.json")	
	if err != nil {
		return err
	}

	defer defaultSrc.Close()

	target, err := os.Create("./jablkoconfig.json")
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, defaultSrc)
	if err != nil {
		return err
	}

	return nil
}