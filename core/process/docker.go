package process

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
)

func dockerImageExists(name string, tag string) bool {
	cmd := exec.Command("docker", "inspect", name+":"+tag)
	err := cmd.Run()

	if err != nil {
		return false
	}

	return true
}

func pullDockerImage(name string, tag string) error {
	cmd := exec.Command("docker", "pull",
		name+":"+tag)

	return cmd.Run()
}

func CreateDockerProcess(name string, tag string) *DockerProcess {
	return &DockerProcess{sync.RWMutex{}, name, tag, 0, nil}
}

// Only pull docker image if the image is not found locally
// This function is run everytime a Docker Seed is started
func imageSafeInstall(name string, tag string) error {
	exists := dockerImageExists(name, tag)

	if exists {
		return nil
	}

	err := pullDockerImage(name, tag)
	if err != nil {
		return fmt.Errorf("Unable to pull Docker image \"%s\"", name+":"+tag)
	}

	return nil
}

func (proc *DockerProcess) Start(port int) error {
	proc.Lock()
	defer proc.Unlock()

	// Check if image is already started

	// Check if docker image doesnt exist and pull
	err := imageSafeInstall(proc.name, proc.tag)
	if err != nil {
		return err
	}

	// Create Cmd
	proc.Cmd = exec.Command("docker", "run",
		"-p", strconv.Itoa(port)+":8080",
		proc.name+":"+proc.tag)

	err = proc.Cmd.Start()
	go func() {
		err := proc.Cmd.Wait()
		log.Info().
			Err(err).
			Msg("Docker process exited")
	}()

	return err
}

func (proc *DockerProcess) Stop() error {
	if proc.Cmd == nil || proc.Cmd.Process == nil {
		return nil
	}

	return proc.Cmd.Process.Signal(os.Interrupt)
}

func (proc *DockerProcess) Name() string {
	return proc.name
}

func (proc *DockerProcess) Tag() string {
	proc.RLock()
	defer proc.RUnlock()
	return proc.tag
}

func (proc *DockerProcess) Type() string {
	return PROC_DOCKER
}

func (proc *DockerProcess) Port() int {
	proc.RLock()
	defer proc.RUnlock()
	return proc.port
}

func (proc *DockerProcess) Update(name string, tag string) error {
	proc.Lock()
	defer proc.Unlock()

	err := proc.Stop()
	if err != nil {
		return err
	}

	return nil
	// TODO
}
