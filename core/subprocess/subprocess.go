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
	Cmd    *exec.Cmd
	Port   int
	Key    string
	Config []byte
	Dir    string
	Env    []string
}

func CreateSubprocess(source string, jablkoPort int, processPort int, jmodKey string, dataDir string, config []byte) *Subprocess {
	// Creates a subprocess from the given parameters
	// Does not start the process

	log.Info().
		Str("subprocess", source).
		Msg("Creating subprocess...")

	sub := new(Subprocess)
	sub.Port = processPort
	sub.Key = jmodKey
	sub.Config = config
	sub.Dir = source
	sub.Env = []string{
		"JABLKO_CORE_PORT=" + strconv.Itoa(jablkoPort),
		"JABLKO_MOD_PORT=" + strconv.Itoa(processPort),
		"JABLKO_MOD_KEY=" + jmodKey,
		"JABLKO_MOD_DATA_DIR=" + dataDir,
		"JABLKO_MOD_CONFIG=" + string(config),
	}

	sub.Cmd = exec.Command("./jablkostart.sh")
	sub.Cmd.Dir = sub.Dir
	sub.Cmd.Env = sub.Env

	sub.Cmd.Stdout = ColoredWriter{os.Stdout}
	sub.Cmd.Stderr = ColoredWriter{os.Stderr}

	return sub
}

func (sub *Subprocess) MarshalJSON() ([]byte, error) {
	return sub.Config, nil
}

func (sub *Subprocess) Start() {
	err := sub.Cmd.Run()

	log.Warn().
		Err(err).
		Int("exitCode", sub.Cmd.ProcessState.ExitCode()).
		Msg("Process exited")

}

func (sub *Subprocess) Restart() {
	sub.Cmd = exec.Command("./jablkostart.sh")
	sub.Cmd.Dir = sub.Dir
	sub.Cmd.Env = sub.Env
	sub.Cmd.Stdout = ColoredWriter{os.Stdout}
	sub.Cmd.Stderr = ColoredWriter{os.Stderr}

	go sub.Start()
}

func (sub *Subprocess) Build() error {
	buildProc := exec.Command("./jablkobuild.sh")
	buildProc.Dir = sub.Cmd.Dir

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
