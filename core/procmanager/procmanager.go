package procmanager

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/ccoverstreet/Jablko/core/process"
)

type ProcManager struct {
	sync.RWMutex
	mods            map[string]process.ModProcess
	portSearchIndex int // Starting value for port search, incremented to prevent races
}

func CreateProcManager() ProcManager {
	return ProcManager{
		sync.RWMutex{},
		make(map[string]process.ModProcess),
		10000,
	}
}

func (pman *ProcManager) MarshalJSON() ([]byte, error) {
	pman.RLock()
	defer pman.RUnlock()
	return json.Marshal(pman.mods)
}

func (pman *ProcManager) UnmarshalJSON(b []byte) error {
	// This struct is used to switch between config types and handlers
	type idStruct struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	rawConfig := map[string]json.RawMessage{}
	idMap := map[string]idStruct{}

	// Unmarshal raw map
	err := json.Unmarshal(b, &rawConfig)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal raw config map: %v", err)
	}

	// Unmarshal id map
	err = json.Unmarshal(b, &idMap)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal id map: %v", err)
	}

	for modName, modID := range idMap {
		fmt.Println(modName, modID, string(rawConfig[modName]))

		switch modID.Type {
		case process.PROC_DOCKER:
			tempProc, err := process.CreateDockerProcessFromBytes(rawConfig[modName])
			if err != nil {
				return err
			}
			pman.mods[modID.Name] = tempProc
		default:
			return fmt.Errorf("Unable to load config. Process type %s is invalid", modID.Type)
		}
	}

	/*
		for name, conf := range tmp {
			switch conf.Type {
			case process.PROC_DEBUG:
				tempProc, _ := process.CreateDebugProcess(name, conf)
				pman.mods[name] = tempProc

			case process.PROC_DOCKER:
				tempProc, _ := process.CreateDockerProcess(name, conf)
				pman.mods[name] = tempProc
			default:
				return fmt.Errorf("Unable to load config. Process type %s is invalid", conf.Type)
			}
		}
	*/

	return nil
}

// TODO: Figure out what type conf should actually be
func (pman *ProcManager) AddMod(procName string, procType string, conf []byte) error {
	pman.Lock()
	defer pman.Unlock()

	if _, ok := pman.mods[procName]; ok {
		return fmt.Errorf("Mod %s already exists", procName)
	}

	var newProc process.ModProcess
	var procCreateErr error

	switch conf.Type {
	case process.PROC_DEBUG:
		newProc, procCreateErr = process.CreateDebugProcessFromBytes(conf)

	case process.PROC_DOCKER:
		newProc, procCreateErr = process.CreateDockerProcessFromBytes(conf)
		/*
			case process.PROC_DOCKER:
				newProc = process.CreateDockerProcess(procName, tag)
		*/
	default:
		return fmt.Errorf("Invalid mod type specified")
	}

	if procCreateErr != nil {
		return procCreateErr
	}

	pman.mods[procName] = newProc

	return nil
}

// Stop the mod, remove from mods
func (pman *ProcManager) RemoveMod(procName string) error {
	pman.Lock()
	defer pman.Unlock()

	proc, ok := pman.mods[procName]

	if !ok {
		return fmt.Errorf("Unable to remove mod. %s does not exist")
	}

	err := proc.Stop()
	if err != nil {
		return err
	}

	delete(pman.mods, procName)
	return nil
}

func (pman *ProcManager) UpdateMod(procName string, tag string) error {
	pman.Lock()
	defer pman.Unlock()

	proc, ok := pman.mods[procName]
	if !ok {
		return fmt.Errorf("Unable to update mod. Mod %s does not exist", procName)
	}

	return proc.Update(procName, tag)
}

func getOpenPort(start int, stop int) (int, error) {
	for i := start; i < stop; i++ {
		conn, err := net.Listen("tcp", fmt.Sprintf(":%d", i))
		if err == nil {
			conn.Close()
			return i, nil
		}

		conn.Close()
	}

	return 0, fmt.Errorf("No open ports found in range %d to %d", start, stop)
}

func (pman *ProcManager) StartMod(procName string) error {
	pman.Lock()
	defer pman.Unlock()

	proc, ok := pman.mods[procName]
	if !ok {
		return fmt.Errorf("Unable to start mod. Mod %s does not exists", procName)
	}

	port, err := getOpenPort(pman.portSearchIndex, 20000)
	if err != nil {
		return err
	}

	pman.portSearchIndex = ((pman.portSearchIndex + 1) % 10000) + 10000
	fmt.Println(pman.portSearchIndex)

	return proc.Start(port)
}

func (pman *ProcManager) StartAllMods() error {
	errorString := ""
	for modName, _ := range pman.mods {
		err := pman.StartMod(modName)
		if err != nil {
			errorString += "\n" + err.Error()
		}
	}

	return fmt.Errorf("%s", errorString)
}

func (pman *ProcManager) StopMod(procName string) error {
	pman.Lock()
	defer pman.Unlock()

	proc, ok := pman.mods[procName]
	if !ok {
		return fmt.Errorf("Unable to start mod. Mod %s does not exists", procName)
	}

	return proc.Stop()
}

func (pman *ProcManager) PassRequest(modName string, w http.ResponseWriter, r *http.Request) error {
	pman.RLock()
	defer pman.RUnlock()

	proc, ok := pman.mods[modName]
	if !ok {
		return fmt.Errorf("Unable to pass request. Mod %s does not exist", modName)
	}

	return proc.PassRequest(w, r)
}

func (pman *ProcManager) GenerateWCScript() (string, error) {
	pman.RLock()
	defer pman.RUnlock()

	errorString := ""

	wcScript := `
const MOD_COMPONENTS = {}

	`

	for modName, mod := range pman.mods {
		wcStr, err := mod.WebComponent(true)
		if err != nil {
			errorString += "\n" + err.Error()
			continue
		}

		wcScript += fmt.Sprintf("\n\nMOD_COMPONENTS[\"%s\"] = %s\n\n", modName, wcStr)
	}

	return wcScript, fmt.Errorf("%s", errorString)
}
