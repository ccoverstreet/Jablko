package process

import (
	"net/http"
	"os/exec"
	"sync"
)

const (
	PROC_DOCKER = "docker"
)

type SeedProcess interface {
	Start(port int) error
	Stop() error
	Name() string
	Tag() string
	Type() string
	Port() int
	Update(name string, tag string) error // Should stop the Seed, pull/update the Seed
	PassRequest(w http.ResponseWriter, r *http.Request) error
}

type DockerProcess struct {
	sync.RWMutex
	name string // Name of Docker image
	tag  string // Tag of docker image (ex. latest)
	port int
	Cmd  *exec.Cmd
}

// Debug process is used to run the Seed during development
// In this case, the only data stored by Jablko is the port
// used for communication
type DebugProcess struct {
	Port int // Port that will be mapped to the container
}
