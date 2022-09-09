package procmanager

import (
	"encoding/json"
	"fmt"
	"net"
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
	tmp := map[string]process.ModProcessConfig{}

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	for name, conf := range tmp {
		switch conf.Type {
		case process.PROC_DEBUG:
			pman.mods[name] = process.CreateDebugProcess(name, conf)
		default:
			return fmt.Errorf("Unable to load config. Process type %s is invalid", conf.Type)
		}
	}

	return nil
}

// TODO: Figure out what type conf should actually be
func (pman *ProcManager) AddMod(procName string, conf process.ModProcessConfig) error {
	pman.Lock()
	defer pman.Unlock()

	var newProc process.ModProcess

	if _, ok := pman.mods[procName]; ok {
		return fmt.Errorf("Mod %s already exists")
	}

	switch conf.Type {
	case process.PROC_DEBUG:
		newProc = process.CreateDebugProcess(procName, conf)
		/*
			case process.PROC_DOCKER:
				newProc = process.CreateDockerProcess(procName, tag)
		*/
	default:
		return fmt.Errorf("Invalid mod type specified")
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
