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
	Writer   *SubprocessWriter
}

func CreateSubprocess(source string, jablkoPort int, jmodKey string, dataDir string, config []byte) (*Subprocess, error) {
	// Creates a subprocess from the given parameters
	// Does not start the process

	log.Info().
		Str("subprocess", source).
		Msg("Creating subprocess...")

	// Make data directory
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, err
	}

	sub := new(Subprocess)
	sub.CorePort = jablkoPort
	sub.Key = jmodKey
	sub.Config = config
	if sub.Config == nil {
		sub.Config = []byte("{}")
	}
	sub.Dir = source
	sub.DataDir = dataDir

	newWriter, err := CreateSubprocessWriter(source)
	if err != nil {
		return nil, err
	}

	sub.Writer = newWriter

	return sub, nil
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

	hostEnv := os.Environ()
	jablkoEnv := []string{
		"JABLKO_CORE_PORT=" + strconv.Itoa(sub.CorePort),
		"JABLKO_MOD_PORT=" + strconv.Itoa(sub.ModPort),
		"JABLKO_MOD_KEY=" + sub.Key,
		"JABLKO_MOD_DATA_DIR=" + sub.DataDir,
		"JABLKO_MOD_CONFIG=" + string(sub.Config),
	}

	sub.Cmd.Env = append(hostEnv, jablkoEnv...)

	sub.Cmd.Stdout = sub.Writer
	sub.Cmd.Stderr = sub.Writer

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
	portNumber, err := GetAvailablePort(10000, 30000)
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
		Int("port", sub.ModPort).
		Msg("Subprocess starting")

	go sub.wait()

	return nil
}

type ReservedPortMap struct {
	sync.Mutex
	Ports map[int]bool
}

var ReservedPorts = &ReservedPortMap{sync.Mutex{}, make(map[int]bool)}

func GetAvailablePort(minPort int, maxPort int) (int, error) {
	ReservedPorts.Lock()
	defer ReservedPorts.Unlock()

	for i := minPort; i < maxPort; i++ {
		// If port was already reserved by Jablko Process
		if ReservedPorts.Ports[i] {
			continue
		}

		conn, err := net.Listen("tcp", fmt.Sprintf(":%d", i))
		if err == nil {
			conn.Close()
			ReservedPorts.Ports[i] = true
			fmt.Println(ReservedPorts)
			return i, nil
		}
	}

	return 0, fmt.Errorf("Unable to find available port in range")
}

func (sub *Subprocess) wait() {
	err := sub.Cmd.Wait()

	log.Warn().
		Err(err).
		Str("jmodName", sub.Dir).
		Int("exitCode", sub.Cmd.ProcessState.ExitCode()).
		Msg("Process exited")
}

func (sub *Subprocess) Stop() error {
	sub.Lock()
	defer sub.Unlock()

	// If process was never succesfully created
	if sub.Cmd.Process == nil {
		return fmt.Errorf("Process doesn't exist")
	}

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

func (sub *Subprocess) GetCurLogBytes() ([]byte, error) {
	sub.Lock()
	defer sub.Unlock()

	return sub.Writer.GetCurLogBytes()
}
