package modmanager

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ccoverstreet/Jablko/core/process"
	"github.com/rs/zerolog/log"
)

type ModManager struct {
	sync.RWMutex
	Mods map[string]*process.DockerProc
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

func (mm *ModManager) StartJMOD(modName string, port int) error {
	mm.RLock()
	defer mm.RUnlock()

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

//func (mm *ModManager) PassRequestToJMOD(w http.ResponseWriter)
