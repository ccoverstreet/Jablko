package process

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

func CreateDebugProcess(name string, conf ModProcessConfig) *DebugProcess {
	return &DebugProcess{sync.RWMutex{}, name, conf.Tag, conf.Port, ""}
}

func CreateDebugProcessFromBytes(name string, conf []byte) (*DebugProcess, error) {
	var tmp ModProcessConfig

	err := json.Unmarshal(conf, &tmp)

	if err != nil {
		return nil, err
	}

	return CreateDebugProcess(name, tmp), err
}

func (proc *DebugProcess) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{
		"tag": "%s",
		"type": "%s",
		"port": %d
	}`, proc.tag, PROC_DEBUG, proc.port)), nil
}

func (proc *DebugProcess) Start(port int) error {
	return nil
}

func (proc *DebugProcess) Stop() error {
	return nil
}

func (proc *DebugProcess) Update(name string, tag string) error {
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

func (proc *DebugProcess) WebComponent() (string, error) {
	return proc.webComponent, nil
}
