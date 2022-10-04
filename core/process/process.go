package process

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
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
	WebComponent(bool) (string, error)
	MarshalJSON() ([]byte, error)
}

type ModProcessConfig struct {
	Tag  string `json:"tag"`
	Type string `json:"type"`
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

func ProxyHTTPRequest(w http.ResponseWriter, r *http.Request, port int) error {
	url, _ := url.Parse("http://127.0.0.1:" + strconv.Itoa(port))
	proxy := httputil.NewSingleHostReverseProxy(url)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	proxy.ServeHTTP(w, r)

	return nil
}
