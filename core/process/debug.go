package process

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
)

type DebugProcessConf struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
	Port int    `json:"port"`
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

func CreateDebugProcess(conf DebugProcessConf) (*DebugProcess, error) {
	if conf.Port == 0 || conf.Port > 65535 {
		return nil, fmt.Errorf("Invalid port specified for debug-type mod: %d", conf.Port)
	}

	return &DebugProcess{sync.RWMutex{}, conf.Name, conf.Tag, conf.Port, ""}, nil
}

func CreateDebugProcessFromBytes(b []byte) (*DebugProcess, error) {
	conf := DebugProcessConf{}

	err := json.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}

	return CreateDebugProcess(conf)
}

func (proc *DebugProcess) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{
		"tag": "%s",
		"type": "%s",
		"port": %d
	}`, proc.tag, PROC_DEBUG, proc.port)), nil
}

func (proc *DebugProcess) Start(port int) error {
	log.Debug().
		Str("name", proc.name).
		Msg("Starting debug process")

	return nil
}

func (proc *DebugProcess) Stop() error {
	log.Debug().
		Str("name", proc.name).
		Msg("Stopping debug process")

	return nil
}

func (proc *DebugProcess) Update(name string, tag string) error {
	proc.Lock()
	defer proc.Unlock()
	log.Debug().
		Str("name", proc.name).
		Msg("Updating debug process (just changing tag)")

	proc.tag = tag

	return nil
}

func (proc *DebugProcess) Name() string {
	return proc.name
}

func (proc *DebugProcess) Tag() string {
	return PROC_DEBUG
}

func (proc *DebugProcess) Type() string {
	return PROC_DEBUG
}

func (proc *DebugProcess) Port() int {
	return proc.port
}

func (proc *DebugProcess) PassRequest(w http.ResponseWriter, r *http.Request) error {
	proc.RLock()
	defer proc.RUnlock()

	return ProxyHTTPRequest(w, r, proc.port)
}

func (proc *DebugProcess) WebComponent(refresh bool) (string, error) {
	// Update existing webcomponent store
	if refresh {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/webComponent", proc.port))
		if err != nil {
			return "", err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		proc.webComponent = string(body)
	}

	return proc.webComponent, nil
}
