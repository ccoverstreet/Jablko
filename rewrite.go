package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ccoverstreet/Jablko/core/process"
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

	// Temporary test code
	proc, err := process.CreateProc(process.ProcConfig{
		"asd",
		"latest",
	})

	if err != nil {
		panic(err)
	}

	log.Printf("%v", proc.CreateDataDirIfNE())
	proc.Start(8080)

	time.Sleep(10 * time.Second)
	proc.Kill()
}

func setupLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().Logger()
}

func createNeededDirs() error {
	err := os.Mkdir("log", 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir("data", 0755)

	return err
}
