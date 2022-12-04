package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const DEFAULT_CONFIG string = `
{
	"updateInterval": 3
}
`

type TestMod struct {
	UpdateInterval int `json:"updateInterval"`
}

func loadConfig() ([]byte, error) {
	b, err := os.ReadFile("jablko/modconfig.json")
	if err != nil {
		err = os.WriteFile("jablko/modconfig.json", []byte(DEFAULT_CONFIG), 0666)
		if err != nil {
			return nil, err
		}

		return []byte(DEFAULT_CONFIG), nil
	}

	return b, nil
}

func main() {
	conf, err := loadConfig()
	if err != nil {
		panic(err)
	}

	instance := TestMod{}
	err = json.Unmarshal(conf, &instance)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/mod/{func}", demoHandler)
	router.HandleFunc("/webComponent", webComponentHandler)

	go func() {
		for {
			time.Sleep(time.Duration(instance.UpdateInterval) * time.Second)
			log.Println("Doing work")
		}
	}()

	log.Println("Starting HTTP Server")
	http.ListenAndServe(":9090", router)
	fmt.Println("vim-go")
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func demoHandler(w http.ResponseWriter, r *http.Request) {
	fun := mux.Vars(r)["func"]

	res := struct {
		Msg  string `json:"msg"`
		Func string `json:"func"`
	}{
		"Test message received",
		fun,
	}

	log.Printf("Function %s requested", fun)

	sendJSONResponse(w, res)
}

func webComponentHandler(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile("webcomponent.js")
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"err": "Unable to read webcomponent file"}`)
		return
	}

	w.Header().Set("Content-Type", "text/javascript")
	w.Write(b)
}
