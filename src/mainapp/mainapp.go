// mainapp.go: Main application definitions for Jablko
// Cale Overstreet
// 2020/24/11
// Contains main internal definitions for Jablko

package mainapp

import (
	"log"
)

type generalConfig struct {}

var jablkoConfig = generalConfig{HttpPort: 8080, HttpsPort: -1}

type MainApp struct {}

func CreateMainApp(configData []byte) *MainApp {
	log.Println(configData)
}

func (app MainApp) SendMessage(message string) error {
	log.Printf("Message: %s\n", message)	
} 
