// Jablko Module Communication
// Cale Overstreet
// February 2, 2021
// Responsible for registering modules and storing module name, auth string, and ip address

package modcommunication

import (
	"sync"
)

type registryEntry struct {
	address string
	authStr string
}

type modRegistry struct {
	sync.Mutex	
	mods map[string]registryEntry
}

var moduleRegistry modRegistry

func InitializeRegistry() {

}
