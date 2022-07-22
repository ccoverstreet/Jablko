package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ccoverstreet/Jarmuz-Status/app"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Jarmuz Status starting...")

	JablkoCorePort := os.Getenv("JABLKO_CORE_PORT")
	JMODPort := os.Getenv("JABLKO_MOD_PORT")
	JMODKey := os.Getenv("JABLKO_MOD_KEY")
	//JMODDataDir := os.Getenv("JABLKO_MOD_DATA_DIR")
	JMODConfig := os.Getenv("JABLKO_MOD_CONFIG")

	jarmuzstatus := app.CreateStatusApp(JMODConfig, JMODPort, JMODKey, JablkoCorePort)

	// Application Routes
	router := mux.NewRouter()
	router.HandleFunc("/webComponent", app.WebComponentHandler)
	router.HandleFunc("/instanceData", app.InstanceDataHandler)
	router.HandleFunc("/jmod/clientWebsocket", jarmuzstatus.HandleClientWebsocket)
	router.HandleFunc("/jmod/removeDevice", app.WrapHandler(app.RemoveDeviceHandler, jarmuzstatus))
	router.HandleFunc("/jmod/addDevice", app.WrapHandler(app.AddDeviceHandler, jarmuzstatus))

	jarmuzstatus.Run()
	http.ListenAndServe(":"+JMODPort, router)
}
