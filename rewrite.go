package main

import (
	"log"
	"time"

	"github.com/ccoverstreet/Jablko/core/process"
)

func main() {
	proc, err := process.CreateProc(process.ProcConfig{
		"asd",
		"latest",
	})

	if err != nil {
		panic(err)
	}

	log.Println(proc.CreateDataDirIfNE())
	proc.Start(8080)

	time.Sleep(10 * time.Second)
	proc.Kill()

	log.Println(proc)
}
