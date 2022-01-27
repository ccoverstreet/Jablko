// Jablko Proc: Docker process abstraction layer
// Cale Overstreet
// January 23, 2022
// This package serves as the itermediate layer between Docker and Jablko.
// This is responsible for starting the docker process, listening to output,
// and providing the common interface for the rest of Jablko.

package process

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

const JMODSTARTMESSAGE = `
==================== JMOD STARTED ====================
`

const JMODSTOPMESSAGE = `
==================== JMOD STOPPED ====================
`

const JMODKILLEDMESSAGE = `
==================== JMOD KILLED ====================
`

func CleanImageName(imageName string) string {
	return strings.ReplaceAll(imageName, "/", "_")
}

type ProcConfig struct {
	ImageName string
	Tag       string
}

type DockerProc struct {
	sync.RWMutex
	Config  ProcConfig
	DataDir string
	Cmd     *exec.Cmd
	Writer  *JMODWriter
}

func CreateProc(config ProcConfig) (*DockerProc, error) {
	absPath, err := filepath.Abs("data/" + CleanImageName(config.ImageName))
	if err != nil {
		return nil, err
	}

	writer, err := CreateJMODWriter(config.ImageName)
	if err != nil {
		return nil, err
	}

	return &DockerProc{sync.RWMutex{}, config, absPath, nil, writer}, nil
}

func (proc *DockerProc) MarshalJSON() ([]byte, error) {
	proc.Lock()
	defer proc.Unlock()

	var tempStruct = struct {
		Tag string `json:"tag"`
	}{proc.Config.Tag}

	return json.Marshal(tempStruct)
}

func (proc *DockerProc) IsLocal() bool {
	proc.RLock()
	defer proc.RUnlock()

	return proc.Config.Tag == "local"
}

func (proc *DockerProc) PullImage() error {
	log.Info().
		Str("imageName", proc.Config.ImageName).
		Msg("Pulling JMOD image from DockerHub")
	cmd := exec.Command("docker", "pull", proc.Config.ImageName)
	return cmd.Run()
}

func (proc *DockerProc) Start(port int) error {
	proc.Lock()
	defer proc.Unlock()

	proc.CreateDataDirIfNE()
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

	fmt.Fprintf(proc.Writer, JMODSTARTMESSAGE)

	go proc.wait()

	return nil
}

func (proc *DockerProc) wait() {
	proc.Cmd.Wait()
	fmt.Fprintf(proc.Writer, JMODSTOPMESSAGE)
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

	fmt.Fprintf(proc.Writer, JMODKILLEDMESSAGE)

	if proc.Cmd == nil || proc.Cmd.Process == nil {
		return nil
	}

	return proc.Cmd.Process.Signal(os.Interrupt)
}

func (proc *DockerProc) CreateDataDirIfNE() error {
	_, err := os.Stat(proc.DataDir)

	if err == nil || !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(proc.DataDir, 0755)
}

// Overwrites the config file located within the JMOD image's data directory
func (proc *DockerProc) UpdateConfig(newConfig string) error {

	return nil
}
