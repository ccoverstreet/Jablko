package process

import (
	"net/http"
	"os/exec"
	"sync"
)

const (
	PROC_DEBUG  = "debug"
	PROC_DOCKER = "docker"
)

type ModProcess interface {
	Start(port int) error
	Stop() error
	Update(name string, tag string) error // Should stop the mod, pull/update the mod, and restart
	Name() string
	Tag() string
	Type() string
	Port() int
	PassRequest(w http.ResponseWriter, r *http.Request) error
	WebComponent() (string, error)
	MarshalJSON() ([]byte, error)
}

type ModProcessConfig struct {
	Tag  string `json:"tag"`
	Type string `json:"type"`
	Port int    `json:"port"`
}

type DockerProcess struct {
	sync.RWMutex
	name string // Name of Docker image
	tag  string // Tag of docker image (ex. latest)
	port int
	Cmd  *exec.Cmd
}

// In this case, the only data stored by Jablko is the port
// used for communication
type DebugProcess struct {
	sync.RWMutex
	name         string
	tag          string
	port         int
	webComponent string
}
