// Mod Manager
// Cale Overstreet
// Apr. 24, 2021

// Response for process management for JMODs, passing data to
// JMODs, installing/upgrading JMODS.

package modmanager

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"bytes"

	"github.com/rs/zerolog/log"
	"github.com/buger/jsonparser"

	"github.com/ccoverstreet/Jablko/core/subprocess"
	"github.com/ccoverstreet/Jablko/core/jutil"
)

type ModManager struct {
	ConfigMap map[string][]byte
	ProcMap map[string]*subprocess.Subprocess
}

var curPort = 8081

func NewModManager(conf []byte) (*ModManager, error) {
	newMM := new(ModManager)
	newMM.ConfigMap = make(map[string][]byte)
	newMM.ProcMap = make(map[string]*subprocess.Subprocess)

	parseConfObj := func(key []byte, value []byte, _ jsonparser.ValueType, _ int) error {
		newMM.ConfigMap[string(key)] = value
		return nil
	}

	jsonparser.ObjectEach(conf, parseConfObj)

	// Try to start all subprocesses
	for key, conf := range newMM.ConfigMap {
		jmodKey, err := jutil.RandomString(32)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to generate random string for jmodKey")

			continue
		}

		newSub, err := subprocess.CreateSubprocess(key, 8080, curPort, jmodKey, "./data", conf)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to create subprocess")
			continue
		}
		curPort += 1

		err = newSub.Build()

		err = newSub.Start()
		if err != nil {
			log.Error().
				Err(err).
				Str("subprocess", key).
				Msg("Unable to start subprocess")

			continue
		}

		newMM.ProcMap[key] = newSub
	}

	return newMM, nil;
}

func (mm *ModManager) PassRequest(w http.ResponseWriter, r *http.Request) {
	source := r.FormValue("JMOD-Source")

	modPort := mm.ProcMap[source].Port
	url, _ := url.Parse("http://localhost:" + strconv.Itoa(modPort))
	proxy := httputil.NewSingleHostReverseProxy(url)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to read incoming proxy request body")
	}

	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	proxy.ServeHTTP(w, r)
}
