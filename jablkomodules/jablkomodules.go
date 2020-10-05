// jablkomodules.go: Jablko Module Manager

package jablkomodules

import (
	"fmt"
	"net/http"
)

var HandleMap = make(map[string]func(http.ResponseWriter, *http.Request)(error), 0)
var InitializationMap = make(map[string]func([]byte)(error))
InitializationMap["Jablko-Interface-Status-Go"] = 

func Initialize(jablkoModulesConfig []byte) {
	fmt.Printf("%s\n", jablkoModulesConfig)	
}
