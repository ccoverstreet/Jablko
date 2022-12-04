package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ccoverstreet/Jablko/core/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	greeting("0.3.0")

	// Setup environment
	setupLogging()
	setupDirs()

	jsonstr, err := os.ReadFile("jablkoconfig.json")
	if err != nil {
		panic(err)
	}

	app, err := core.CreateJablkoCore(jsonstr)
	if err != nil {
		panic(err)
	}

	/*
		app.PMan.AddMod("testasjfsadfjhasd", process.ModProcessConfig{
			"sometag",
			process.PROC_DEBUG,
			8080,
	*/

	b, err := json.MarshalIndent(app, "", "    ")
	fmt.Println(err, string(b))
	app.Start()

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

func setupDirs() {
	if _, err := os.Stat("data"); err != nil {
		err := os.MkdirAll("data", 0700)
		if err != nil {
			log.Fatal().
				Err(err).
				Msgf("Cannot create required directory %s", "data")
		}
	}
}
