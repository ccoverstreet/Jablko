// Subprocess Struct and Operations for Jablko
// Cale Overstreet
// Apr. 9, 2021

// The subprocess struct spawned for each Jablko
// Mod. Provides the Jablko-specific environment
// variables

package subprocess

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

type Subprocess struct {
	sync.Mutex
	Cmd      *exec.Cmd
	ModPort  int
	CorePort int
	Key      string
	Config   []byte
	Dir      string
	DataDir  string
}

func CreateSubprocess(source string, jablkoPort int, jmodKey string, dataDir string, config []byte) *Subprocess {
	// Creates a subprocess from the given parameters
	// Does not start the process

	log.Info().
		Str("subprocess", source).
		Msg("Creating subprocess...")

	sub := new(Subprocess)
	sub.CorePort = jablkoPort
	sub.Key = jmodKey
	sub.Config = config
	sub.Dir = source
	sub.DataDir = dataDir

	return sub
}

func (sub *Subprocess) MarshalJSON() ([]byte, error) {
	return sub.Config, nil
}

// Copies parameters stored in the subprocess into a
// new exec.Cmd and sets the environment
// ONLY called when Subprocess.Start is called
func (sub *Subprocess) generateCMD() {
	sub.Cmd = exec.Command("make", "run")
	sub.Cmd.Dir = sub.Dir
	sub.Cmd.Env = []string{
		"JABLKO_CORE_PORT=" + strconv.Itoa(sub.CorePort),
		"JABLKO_MOD_PORT=" + strconv.Itoa(sub.ModPort),
		"JABLKO_MOD_KEY=" + sub.Key,
		"JABLKO_MOD_DATA_DIR=" + sub.DataDir,
		"JABLKO_MOD_CONFIG=" + string(sub.Config),
	}
	sub.Cmd.Stdout = ColoredWriter{sub.Dir, os.Stdout}
	sub.Cmd.Stderr = ColoredWriter{sub.Dir, os.Stderr}

	sub.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func (sub *Subprocess) Start() error {
	sub.Lock()
	defer sub.Unlock()

	// Check if process is already started
	// Safety check just in case sub.Cmd is nil
	if sub.Cmd != nil {
		if sub.Cmd.Process != nil && sub.Cmd.ProcessState == nil {
			return fmt.Errorf("Process is already started")
		}
	}

	// Search for available port
	portNumber, err := findAvailablePort(10000, 30000)
	if err != nil {
		return err
	}

	sub.ModPort = portNumber

	// Generate the exec.Cmd used for the process
	sub.generateCMD()

	err = sub.Cmd.Start()

	if err != nil {
		return err
	}

	log.Info().
		Str("source", sub.Dir).
		Msg("Subprocess starting")

	go sub.wait()

	return nil
}

func findAvailablePort(minPort int, maxPort int) (int, error) {
	for i := minPort; i < maxPort; i++ {
		conn, err := net.Listen("tcp", fmt.Sprintf(":%d", i))
		if err == nil {
			conn.Close()
			return i, nil
		}
	}

	return 0, fmt.Errorf("Unable to find available port in range")
}

func (sub *Subprocess) wait() {
	err := sub.Cmd.Wait()

	log.Warn().
		Err(err).
		Int("exitCode", sub.Cmd.ProcessState.ExitCode()).
		Msg("Process exited")
}

func (sub *Subprocess) Stop() error {
	sub.Lock()
	defer sub.Unlock()

	if sub.Cmd.ProcessState != nil {
		return fmt.Errorf("Process already stopped")
	}

	pgid, err := syscall.Getpgid(sub.Cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15)
	}

	return nil //sub.Cmd.Process.Kill()
}

func (sub *Subprocess) Build() error {
	buildProc := exec.Command("make", "build")
	buildProc.Dir = sub.Dir

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
