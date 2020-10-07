// jpm.go: Jablko Package Manager
// Cale Overstreet
// October 5, 2020
package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/buger/jsonparser"
)

func main() {
	fmt.Printf("JPM\n")	

	jablkoConfig, err := ioutil.ReadFile("../jablko_config.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", jablkoConfig)

	template, err := ioutil.ReadFile("../jablkomodules/jablkomodules.template")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", template)
	
	jablkoModulesConfig, _, _, err := jsonparser.Get(jablkoConfig, "jablkoModules")
	if err != nil {
		panic(err)
	}

	var importString string = "import (\n"
	var handlerMapString string = "var WebHandlerMap = map[string]func(w http.ResponseWriter, r *http.Request) {"

	var packageNumber int = 0
	jsonparser.ObjectEach(jablkoModulesConfig, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		// Iterate through each Jablko Package

		// Split key string to namespace import
		splitKey := strings.Split(string(key), "/")

		formattedName := splitKey[len(splitKey) - 2] + "_" + splitKey[len(splitKey) - 1]

		fmt.Printf("%v\n", splitKey)
				
		// Add package to importString
		importString += "\t" + formattedName + " \"" + string(key) + "\"\n"

		// Add package map to handler map
		if packageNumber == 0 {
			packageNumber += 1
			handlerMapString += "\"" + formattedName + "\": " + formattedName + ".WebHandler"
		} else {
			handlerMapString += ", \"" + formattedName + "\": " + formattedName + ".WebHandler"
		}

		// Get the slice for the Jablko Package
		jablkoPackage, _, _, err := jsonparser.Get(jablkoModulesConfig, string(key))
		if err != nil {
			panic(err)
		}

		// Iterate through all instances of a package
		jsonparser.ObjectEach(jablkoPackage, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			return nil
		})
		return nil
	})

	// Close the strings
	importString += ")\n"
	handlerMapString += "}\n"

	fmt.Printf("%s\n", importString)
	fmt.Printf("%s\n", handlerMapString)

	templateStr := string(template)
	templateStr = strings.Replace(templateStr, "$JABLKO_IMPORTS", importString, 1)
	templateStr = strings.Replace(templateStr, "$HANDLER_MAP", handlerMapString, 1)

	fmt.Printf("%s\n", templateStr)

	err = ioutil.WriteFile("./driver.go", []byte(templateStr), 0744)
	if err != nil {
		panic(err)
	}
}
