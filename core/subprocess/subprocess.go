// Subprocess Struct and Operations for Jablko
// Cale Overstreet
// Apr. 9, 2021

// The subprocess struct spawned for each Jablko
// Mod. Provides the Jablko-specific environment
// variables 

package subprocess

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/rs/zerolog/log"
)


type Subprocess struct {
	Cmd *exec.Cmd
	Port int
}

func CreateSubprocess(source string, jablkoPort int, processPort int, dataDir string, config []byte) (*Subprocess, error) {
		// Creates a subprocess from the given parameters
		// Does not start the process

		log.Info().
			Str("subprocess", source).
			Msg("Creating subprocess...")

		sub := new(Subprocess)
		sub.Cmd = exec.Command("./jablkostart.sh")
		sub.Cmd.Dir = source
		sub.Cmd.Env = []string{
			"JABLKO_CORE_PORT=" + strconv.Itoa(jablkoPort),
			"JABLKO_MOD_PORT=" + strconv.Itoa(processPort),
			"JABLKO_MOD_DATA_DIR=" + dataDir,
			"JABLKO_MOD_CONFIG=" + string(config),
		}

		sub.Cmd.Stdout = ColoredWriter{os.Stdout}
		sub.Cmd.Stderr = ColoredWriter{os.Stderr}

		sub.Port = processPort

		return sub, nil
	}

	func (sub *Subprocess) Start() error {
		err := sub.Cmd.Start()

		return err
	}

	func (sub *Subprocess) Build() error {
		buildProc := exec.Command("./jablkobuild.sh")
		buildProc.Dir = sub.Cmd.Dir

		//err := buildProc.Run()
		out, err := buildProc.CombinedOutput()
		log.Debug().
		Err(err).
		Str("buildOutput", string(out)).
		Msg("Build finished")

		return err
	}

	func (sub *Subprocess) Update() {
		// Handles updating the source
		// Will first look for precompiled options
		// and then resort to a build from source
	}
