package main

import (
	"fmt"
	"time"

	"github.com/ccoverstreet/Jablko/core/process"
)

func main() {
	fmt.Println("vim-go")

	x := process.CreateDockerProcess("ccoverstreet/go-sample", "latest")
	fmt.Println(x)

	fmt.Println(x.Start(10000))

	time.Sleep(10 * time.Second)
	fmt.Println(x.Stop())
	time.Sleep(1 * time.Second)
}
