package process

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
)

func CreateDebugProcess(name string, conf ModProcessConfig) (*DebugProcess, error) {
	if conf.Port == 0 || conf.Port > 65535 {
		return nil, fmt.Errorf("Invalid port specified for debug-type mod: %d", conf.Port)
	}
	return &DebugProcess{sync.RWMutex{}, name, conf.Tag, conf.Port, ""}, nil
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

	url, _ := url.Parse("http://127.0.0.1:" + strconv.Itoa(proc.port))
	proxy := httputil.NewSingleHostReverseProxy(url)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	proxy.ServeHTTP(w, r)

	return nil
}

func (proc *DebugProcess) WebComponent(refresh bool) (string, error) {
	// Update existing webcomponent store
	if refresh {
		resp, err := http.Get("http://127.0.0.1:" + strconv.Itoa(proc.port))
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
