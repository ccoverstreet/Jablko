package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ccoverstreet/Jablko/core/process"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	greeting("0.3.0")
	setupLogging()

	x := process.CreateDockerProcess("ccoverstreet/go-sample", "asd")
	fmt.Println(x)

	log.Printf("%v", x.Start(10000))

	time.Sleep(1 * time.Second)
	fmt.Println(x.Stop())
	time.Sleep(1 * time.Second)
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
