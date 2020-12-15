// mainapp.go: Main application definitions for Jablko
// Cale Overstreet
// 2020/24/11
// Contains main internal definitions for Jablko

package mainapp

import (
	"log"

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
}

func CreateMainApp(configData []byte) (*MainApp, error) {
	instance := new(MainApp)

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
	if err != nil {
		panic(err)
	}

	instance.ModHolder = newModHolder

	instance.Db = database.Initialize()

	return instance, nil
}

func (app *MainApp) SendMessage(message string) error {
	log.Printf("Message: %s\n", message)	

	return nil
} 

func (app *MainApp) SyncConfig(modName string) {
	log.Printf("Sync config called for module \"%s\"\n", modName)		
}
