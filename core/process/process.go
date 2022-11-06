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

// Determines proc type from a byte version of the provided config
func DetermineProcType(conf []byte) (string, error) {
	temp := struct {
		Type string `json:"type"`
	}{""}

	err := json.Unmarshal(conf, &temp)

	if temp.Type != PROC_DEBUG && temp.Type != PROC_DOCKER {
		return "", fmt.Errorf("Invalid mod type '%s' specified", temp.Type)
	}

	return temp.Type, err
}

func DetermineProcName(conf []byte) string {
	temp := struct {
		Name string `json:"name"`
	}{}

	json.Unmarshal(conf, &temp)

	return temp.Name
}

func DetermineProcTag(conf []byte) string {
	temp := struct {
		Tag string `json:"tag"`
	}{}

	json.Unmarshal(conf, &temp)

	return temp.Tag
}
