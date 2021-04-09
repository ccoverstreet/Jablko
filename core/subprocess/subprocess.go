// Subprocess Struct and Operations for Jablko
// Cale Overstreet
// Apr. 9, 2021

// The subprocess struct spawned for each Jablko
// Mod. Provides the Jablko-specific environment
// variables 

package subprocess

import (
	"os/exec"
)

type Subprocess {
	Cmd exec.Cmd	
	Port int
}

func CreateSubprocess(source string, port int) error {	
	

	return nil
}

func (proc *Subprocess) Update() {
	// Handles updating the source
	// Will first look for precompiled options
	// and then resort to a build from source
}
