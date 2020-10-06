// jablkomodules.go: Jablko Module Manager

package jablkomodules

import (
	"fmt"
	"net/http"
)

import (
	ccoverstreet_jablkointerfacestatus "github.com/ccoverstreet/jablkointerfacestatus"
)


var HandlerMap = map[string]func(w http.ResponseWriter, r *http.Request) {"ccoverstreet_jablkointerfacestatus": ccoverstreet_jablkointerfacestatus.RouteHandler}


func Initialize(jablkoModulesConfig []byte) {
	fmt.Printf("%s\n", jablkoModulesConfig)	
}
