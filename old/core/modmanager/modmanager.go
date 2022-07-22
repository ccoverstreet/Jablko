package modmanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ccoverstreet/Jablko/core/process"
	"github.com/rs/zerolog/log"
)

type ModManager struct {
	sync.RWMutex
	Mods       map[string]*process.DockerProc
	fRegenDash bool
	dash       string
}

func (mm *ModManager) UnmarshalJSON(data []byte) error {
	tempMap := make(map[string]struct {
		Tag string `json:"tag"`
	})

	err := json.Unmarshal(data, &tempMap)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	mm.Mods = make(map[string]*process.DockerProc)

	errs := []error{}

	for modName, modConf := range tempMap {
		procConf := process.ProcConfig{modName, modConf.Tag}
		proc, err := process.CreateProc(procConf)
		if err != nil {
			errs = append(errs, err)
		}
		mm.Mods[modName] = proc
	}

	mm.fRegenDash = true
	mm.dash = ""

	return nil
}

func (mm *ModManager) MarshalJSON() ([]byte, error) {
	return json.Marshal(mm.Mods)
}

func (mm *ModManager) Cleanup() error {
	for _, mod := range mm.Mods {
		err := mod.Kill()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to stop JMOD during cleanup")
		}
	}

	return nil
}

func findAvailablePort() (int, error) {
	minPort := 10000
	maxPort := 20000

	for i := minPort; i < maxPort; i++ {
		conn, err := net.Listen("tcp", fmt.Sprintf(":%d", i))
		if err == nil {
			conn.Close()
			return i, nil
		}
	}

	return 0, fmt.Errorf("No available port found within range")
}

func (mm *ModManager) StartJMOD(modName string, port int) error {
	mm.RLock()
	defer mm.RUnlock()

	if port == 0 {
		var err error
		port, err = findAvailablePort()
		if err != nil {
			return err
		}
	}

	log.Info().Str("imageName", modName).Msg("Starting JMOD.")

	mod, ok := mm.Mods[modName]
	if !ok {
		return fmt.Errorf(`JMOD with name "%s" does not exist.`, modName)
	}

	if !mod.IsLocal() {
		err := mod.PullImage()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to pull image for JMOD.")
		}
	}

	err := mod.Start(port)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error starting JMOD.")
	}

	return nil
}

func (mm *ModManager) StartAll() {
	mm.RLock()
	defer mm.RUnlock()

	i := 10000

	for modName, _ := range mm.Mods {
		mm.StartJMOD(modName, i)
		i++
	}
}

func (mm *ModManager) StopJMOD(modName string) error {
	mod, ok := mm.Mods[modName]
	if !ok {
		return fmt.Errorf("JMOD '%s' does not exist", modName)
	}

	return mod.Kill()
}

// This should be moved into process
func (mm *ModManager) PassReqToJMOD(w http.ResponseWriter, r *http.Request) error {
	mm.RLock()
	defer mm.RUnlock()

	jmodName := r.FormValue("JMOD")
	proc, ok := mm.Mods[jmodName]
	if !ok {
		return fmt.Errorf("JMOD does not exist.")
	}

	return proc.PassRequest(w, r)
}

func CleanModName(name string) string {
	replacer := strings.NewReplacer(
		"_", "-",
		"/", "-",
	)

	return replacer.Replace(name)
}

// TODO: Properly wrap collected dashboard code
func (mm *ModManager) GenerateDashboard() error {
	jablkoRoot, err := os.Executable()
	if err != nil {
		return err
	}

	fmt.Println(jablkoRoot)

	// Load dashboard template
	dashTemplate, err := ioutil.ReadFile("html/dashboard_template.html")
	if err != nil {
		return err
	}

	formatStr := "WCMAP[\"MODNAME\"] = WEBCOMPONENT\n\n"
	wcs := ""
	errors := ""
	for modName, mod := range mm.Mods {
		wcText, err := mod.WebComponent()
		if err != nil {
			errors += fmt.Sprintf(" %v;", err)
		}

		wcReplacer := strings.NewReplacer(
			"MODNAME", CleanModName(modName),
			"WEBCOMPONENT", wcText,
		)

		wcs += wcReplacer.Replace(formatStr)
	}

	fmt.Println(wcs)

	if len(errors) != 0 {
		return fmt.Errorf(errors)
	}

	mm.dash = strings.Replace(string(dashTemplate), "$WEBCOMPONENTS", wcs, 1)

	return nil
}

func (mm *ModManager) GetDashboard() string {
	mm.RLock()
	defer mm.RUnlock()

	if mm.fRegenDash {
		log.Info().
			Msg("Generating dashboard")

		err := mm.GenerateDashboard()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error generating dashboard")
		}
	}

	return mm.dash
}
