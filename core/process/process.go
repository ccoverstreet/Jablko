// Jablko Proc: Docker process abstraction layer
// Cale Overstreet
// January 23, 2022
// This package serves as the itermediate layer between Docker and Jablko.
// This is responsible for starting the docker process, listening to output,
// and providing the common interface for the rest of Jablko.

package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

func CleanImageName(imageName string) string {
	return strings.ReplaceAll(imageName, "/", "_")
}

type ProcConfig struct {
	ImageName string
	Tag       string
}

type Proc interface {
	Start(port int) error // Starts the JMOD after searching for an available port
	Kill() error
	CreateDataDirIfNE() error
}

type DockerProc struct {
	sync.Mutex
	Config  ProcConfig
	DataDir string
	Cmd     *exec.Cmd
	Writer  *JMODWriter
}

func CreateProc(config ProcConfig) (Proc, error) {
	absPath, err := filepath.Abs(CleanImageName(config.ImageName))
	if err != nil {
		return nil, err
	}

	writer, err := CreateJMODWriter(config.ImageName)
	if err != nil {
		return nil, err
	}

	return &DockerProc{sync.Mutex{}, config, absPath, nil, writer}, nil
}

func (proc *DockerProc) Start(port int) error {
	log.Info().
		Int("port", port).
		Msg("Starting process")

	proc.Cmd = exec.Command("docker", "run",
		"-p", strconv.Itoa(port)+":8080",
		"--mount", fmt.Sprintf("type=bind,source=%s,target=/data", proc.DataDir),
		proc.Config.ImageName)

	//proc.Cmd = exec.Command("echo", "aSAdasdaSDASD")
	proc.Cmd.Stdout = proc.Writer
	proc.Cmd.Stderr = proc.Writer

	err := proc.Cmd.Start()
	if err != nil {
		return err
	}

	fmt.Fprintf(proc.Writer, `
==================== JMOD STARTED ====================
`)

	go proc.wait()

	return nil
}

func (proc *DockerProc) wait() {
	proc.Cmd.Wait()
	fmt.Fprintf(proc.Writer, `
==================== JMOD STOPPED ====================
`)
	log.Info().
		Str("imageName", proc.Config.ImageName).
		Msg("Docker process exited.")
}

// Does not work on Windows
// TODO: Need to find a work around that works for Windows as well
// Kind of low priority
func (proc *DockerProc) Kill() error {
	proc.Lock()
	defer proc.Unlock()

	fmt.Fprintf(proc.Writer, `
==================== JMOD KILLED ====================
`)

	return proc.Cmd.Process.Signal(os.Interrupt)
}

func (proc *DockerProc) CreateDataDirIfNE() error {
	_, err := os.Stat(proc.DataDir)

	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(proc.DataDir, 0755)
}

// Overwrites the config file located within the JMOD image's data directory
func (proc *DockerProc) UpdateConfig(newConfig string) error {

	return nil
}
