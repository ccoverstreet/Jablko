package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"

	"github.com/ccoverstreet/Jarmuz-Cookbook/cookbook"
)

func main() {
	JablkoCorePort := os.Getenv("JABLKO_CORE_PORT")
	JMODPort := os.Getenv("JABLKO_MOD_PORT")
	JMODKey := os.Getenv("JABLKO_MOD_KEY")
	JMODDataDir := os.Getenv("JABLKO_MOD_DATA_DIR")
	JMODConfig := os.Getenv("JABLKO_MOD_CONFIG")

	if len(JablkoCorePort) == 0 {
		log.Println("ERROR: JABLKO_CORE_PORT not set")
		return
	}

	app := cookbook.CreateCookbook(JablkoCorePort, JMODPort, JMODKey, JMODDataDir, JMODConfig)

	log.Println(http.ListenAndServe(":"+JMODPort, app.GetRouter()))
}
