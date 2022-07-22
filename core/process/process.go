package process

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

type ProcessType int64

const (
	Docker ProcessType = iota
)

type ModProcess interface {
	Start(port int) error
	Stop() error
	Name() string
	Tag() string
	Type() string
	Update(name string, tag string) error // Should stop the JMOD, pull/update the JMOD
}

type DockerProcess struct {
	sync.RWMutex
	Name string // Name of Docker image
	Tag  string // Tag of docker image (ex. latest)
	Cmd  *exec.Cmd
}

// Debug process is used to run the JMOD during development
// In this case, the only data stored by Jablko is the port
// used for communication
type DebugProcess struct {
	Port int // Port that will be mapped to the container
}

func DockerImageExists(name string, tag string) error {
	cmd := exec.Command("docker", "inspect", name+":"+tag)
	return cmd.Run()
}

func PullDockerImage(name string, tag string) error {
	cmd := exec.Command("docker", "pull",
		name+":"+tag)
	return cmd.Run()
}

func CreateDockerProcess(name string, tag string) *DockerProcess {
	return &DockerProcess{sync.RWMutex{}, name, tag, nil}
}

func (proc *DockerProcess) Start(port int) error {
	// Create Cmd
	proc.Cmd = exec.Command("docker", "run",
		"-p", strconv.Itoa(port)+":8080",
		proc.Name+":"+proc.Tag)

	err := proc.Cmd.Start()
	go func() {
		err := proc.Cmd.Wait()
		fmt.Println("Docker process exited:", err)
	}()

	return err
}

func (proc *DockerProcess) Stop() error {
	return proc.Cmd.Process.Signal(os.Interrupt)
}
