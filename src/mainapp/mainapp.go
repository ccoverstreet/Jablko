// mainapp.go: Main application definitions for Jablko
// Cale Overstreet
// 2020/24/11
// Contains main internal definitions for Jablko

package mainapp

import (
	"log"
	"strconv"
	"strings"
	"io/ioutil"
	"os"

	"github.com/buger/jsonparser"

	"github.com/ccoverstreet/Jablko/src/jablkomods"
	"github.com/ccoverstreet/Jablko/src/database"
)

type generalConfig struct {
	HttpPort int
	HttpsPort int
}

var jablkoConfig = generalConfig{HttpPort: 8080, HttpsPort: -1}

type MainApp struct {
	Config generalConfig
	ModHolder *jablkomods.JablkoModuleHolder
	Db *database.JablkoDB
	flags map[string]bool
}

func CreateMainApp(configFilePath string) (*MainApp, error) {

	configData, err := ioutil.ReadFile("./jablkoconfig.json")
	if err != nil {
		log.Printf("%v\n", err)
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
		log.Printf("%v\n", err)
		panic("Error getting HTTP port data\n")
	}

	httpsPort, err := jsonparser.GetInt(configData, "https", "port")
	if err != nil {
		log.Printf("HTTPS port Config not set in Config file\n")
	} else {
		jablkoConfig.HttpsPort = int(httpsPort)
	}

	instance.Config.HttpPort = int(httpPort)
	instance.Config.HttpsPort = int(httpsPort)

	jablkoModulesSlice, _, _, err := jsonparser.Get(configData, "jablkoModules")
	if err != nil {
		panic("Error get Jablko Modules Config\n")
	}

	moduleOrderSlice, _, _, err := jsonparser.Get(configData, "moduleOrder")

	newModHolder, err := jablkomods.Initialize2(jablkoModulesSlice, moduleOrderSlice, instance)

	log.Println(newModHolder)
	if err != nil {
		log.Printf("JablkoMods ERROR: %s\n", err)
	}

	instance.ModHolder = newModHolder

	instance.Db = database.Initialize()

	return instance, nil
}

func (app *MainApp) SendMessage(message string) error {
	log.Printf("Message: %s\n", message)	

	return nil
} 

func (app *MainApp) GetFlagValue(flag string) bool {
	if val, ok := app.flags[flag]; ok {
		return val
	} else {
		return false
	}
}

func (app *MainApp) SyncConfig(modId string) {
	log.Printf("Sync config called for module \"%s\"\n", modId)		
	log.Println("Initial")
	log.Println(app.ModHolder.Config[modId])

	ConfigTemplate:= `{
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

	if _, ok := app.ModHolder.Mods[modId]; !ok {
		log.Printf("Cannot find module %s", modId)
		return 
	}

	newConfByte, err := app.ModHolder.Mods[modId].ConfigStr()
	newConfStr := string(newConfByte)
	if err != nil {
		log.Printf("Unable to get Config string for module %s\n", modId)
	}

	if app.ModHolder.Config[modId] == newConfStr {
		// If there is no change in config
		return 
	}


	app.ModHolder.Config[modId] = newConfStr

	log.Println(string(newConfStr))

	log.Println("Updated")
	log.Println(app.ModHolder.Config)

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

	log.Println(orderStr)

	r := strings.NewReplacer("$httpPort", strconv.Itoa(jablkoConfig.HttpPort),
	"$httpsPort", strconv.Itoa(jablkoConfig.HttpsPort),
	"$moduleString", jablkoModulesStr,
	"$moduleOrder", orderStr)

	ConfigDumpStr := r.Replace(ConfigTemplate)

	err = ioutil.WriteFile("./jablkoconfig.json", []byte(ConfigDumpStr), 0022)
	if err != nil {
		log.Println(`Unable to write to "jablkoconfig.json".`)
	}
}
