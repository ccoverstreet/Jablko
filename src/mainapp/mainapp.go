// mainapp.go: Main application definitions for Jablko
// Cale Overstreet
// 2020/24/11
// Contains main internal definitions for Jablko

package mainapp

import (
	"log"

	"github.com/buger/jsonparser"
)

type generalConfig struct {
	HttpPort int
	HttpsPort int
}

var jablkoConfig = generalConfig{HttpPort: 8080, HttpsPort: -1}

type MainApp struct {
	Config generalConfig
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

	return instance, nil
}

func (app MainApp) SendMessage(message string) error {
	log.Printf("Message: %s\n", message)	

	return nil
} 
