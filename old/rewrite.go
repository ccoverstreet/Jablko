package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/ccoverstreet/Jablko/core/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const startingStatement = `
Jablko Smart Home System
Cale Overstreet
`

func main() {
	fmt.Printf("%s\n", startingStatement)

	// Environment setup code
	setupLogging()
	createNeededDirs()

	confB, err := ioutil.ReadFile("./jablkoconfig.json")
	if err != nil {
		panic(err)
	}

	core, err := app.CreateJablkoCore(confB)
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(core, "", "    ")
	if err != nil {
		panic(err)
	}

	log.Printf("%s", b)

	setupInterruptHandler(core)
	core.StartAllMods()
	core.Listen()
	core.Cleanup()
}

func setupLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().Logger()
}

func createNeededDirs() error {
	err := os.Mkdir("log", 0755)
	err = os.Mkdir("data", 0755)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	return err
}

func setupInterruptHandler(core *app.JablkoCore) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		core.Cleanup()
		os.Exit(0)
	}()
}
