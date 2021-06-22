package main

import (
	"fmt"

	"github.com/ccoverstreet/Jablko/core/github"
)

func main() {
	name := "https://api.github.com/repos/ccoverstreet/Jablko/zipball/refs/tags/v0.2.0"
	/*
		err := github.DownloadZipRepo(name)
		if err != nil {
			panic(err)
		}
	*/

	err := github.UnpackZipRepo(name, "github.com/ccoverstreet/hamstermonitor")
	fmt.Println(err)
}
