package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ccoverstreet/Jablko/core/core"
	"github.com/ccoverstreet/Jablko/core/process"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	greeting("0.3.0")
	setupLogging()

	jsonstr, err := os.ReadFile("jablkoconfig.json")
	if err != nil {
		panic(err)
	}

	app, err := core.CreateJablkoCore(jsonstr)
	if err != nil {
		panic(err)
	}

	app.PMan.AddMod("testasjfsadfjhasd", process.ModProcessConfig{
		"sometag",
		process.PROC_DEBUG,
		8080,
	})

	b, err := json.MarshalIndent(app, "", "    ")
	fmt.Println(err, string(b))
	//app.Start()

	/*
		pman := procmanager.CreateProcManager()
		fmt.Println(pman)
		pman.AddMod(process.PROC_DEBUG, "tester", "")
		pman.AddMod(process.PROC_DEBUG, "teste", "")
		pman.AddMod(process.PROC_DEBUG, "test", "")
		err := pman.StartMod("tester")
		fmt.Println(err)
		pman.StartMod("teste")
		pman.StartMod("test")
		fmt.Println(pman)
		pman.RemoveMod("tester")
		fmt.Println(pman)
	*/
}

func greeting(version string) {
	fmt.Printf(`
	Jablko Core
	Version %s

`, version)
}

func setupLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().Logger()
}
