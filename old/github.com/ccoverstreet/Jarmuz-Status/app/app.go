package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type DeviceInfo struct {
	Name string `json:"name"`
}

type StatusApp struct {
	sync.RWMutex
	PortHTTP      string                `json:"-"`
	Devices       map[string]DeviceInfo `json:"devices"`
	PollInterval  int                   `json:"pollInterval"` // Polling interval in seconds
	status        []bool
	SaveConfig    func([]byte) error `json:"-"`
	connList      []*websocket.Conn
	statusSummary []byte // Stored as JSON
}

var defaultConfig string = `
{
	"pollInterval": 300,
	"devices": {}
}
`

func wrapperSaveConfig(JMODPort, JMODKey, JablkoCorePort string) func([]byte) error {
	return func(config []byte) error {
		client := &http.Client{}
		reqPath := "http://localhost:" + JablkoCorePort + "/service/saveConfig"

		req, err := http.NewRequest("POST", reqPath, bytes.NewBuffer(config))
		if err != nil {
			return err
		}

		req.Header.Add("JMOD-KEY", JMODKey)
		req.Header.Add("JMOD-PORT", JMODPort)

		_, err = client.Do(req)
		log.Println(err)

		return err
	}
}

func CreateStatusApp(config string, JMODPort string, JMODKey string, JablkoCorePort string) *StatusApp {
	app := &StatusApp{}
	app.SaveConfig = wrapperSaveConfig(JMODPort, JMODKey, JablkoCorePort)
	app.PortHTTP = JMODPort

	log.Println(config)

	// If no config is provided
	// Should initialize with default config and save config
	if len(config) < 5 {
		config = defaultConfig
		err := app.SaveConfig([]byte(defaultConfig))
		if err != nil {
			log.Println("Unable to save default config")
			log.Println(err)
			panic(err)
		}
	}

	err := json.Unmarshal([]byte(config), &app)
	if err != nil {
		log.Println("Unable to use provided config")
		log.Println(err)
		panic(err)
	}

	log.Println(app)

	// Default empty status
	app.statusSummary = []byte("[]")

	return app
}

func (app *StatusApp) Save() error {
	newConfig, err := json.Marshal(app)
	if err != nil {
		return err
	}

	return app.SaveConfig(newConfig)
}

func (app *StatusApp) Run() {
	go app.Poll()
	log.Println("Jarmuz Status polling...")
}

func (app *StatusApp) Poll() {
	for {
		app.UpdateSummary()
		app.PushConnections()

		time.Sleep(time.Duration(app.PollInterval) * time.Second)
	}
}

func (app *StatusApp) UpdateSummary() {
	type statusInfo struct {
		IP       string
		Name     string
		IsOnline bool
	}
	outputData := make([]statusInfo, len(app.Devices))

	// --------- Parallel Connection Testing ---------
	wg := sync.WaitGroup{}
	wg.Add(len(app.Devices))
	i := 0
	for devIP, devInfo := range app.Devices {
		go func(id int, statusData statusInfo, wg *sync.WaitGroup) {
			defer wg.Done()
			statusData.IsOnline = Ping(statusData.IP)
			outputData[id] = statusData
		}(i, statusInfo{devIP, devInfo.Name, false}, &wg)
		i = i + 1
	}

	wg.Wait()
	// --------- END Parallel Connection Testing ---------

	summary, err := json.Marshal(outputData)
	if err != nil {
		log.Printf("ERROR: Unable to marshal status output - %v\n", err)
		return
	}

	app.statusSummary = summary
}

func Ping(ipAddress string) bool {
	_, err := net.DialTimeout("tcp", ipAddress+":80", time.Duration(3*time.Second))
	if err != nil && !strings.Contains(err.Error(), "connect: connection refused") {
		return false
	}

	return true
}

var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

func (app *StatusApp) HandleClientWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("Client connecting to live view...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	//defer conn.Close()

	app.addConnection(conn)

	conn.WriteMessage(1, app.statusSummary)

	/*
		for {
			messageType, message, err := conn.ReadMessage()
			log.Println(messageType, message)
			if err != nil {
				log.Println("Error reading client")
				return
			}
		}
	*/
}

func (app *StatusApp) addConnection(conn *websocket.Conn) {
	app.Lock()
	defer app.Unlock()

	app.connList = append(app.connList, conn)
}

// Pushs connection list and status to all connected clients
func (app *StatusApp) PushConnections() {
	app.RLock()
	defer app.RUnlock()

	delMap := make(map[int]struct{})

	for i, conn := range app.connList {
		// Delete connection if nil
		if conn == nil {
			delMap[i] = struct{}{}
			continue
		}
		err := conn.WriteMessage(1, app.statusSummary)
		if err != nil { // {
			delMap[i] = struct{}{}
		}
	}

	app.connList = removeConns(app.connList, delMap)
}

func removeConns(conns []*websocket.Conn, indices map[int]struct{}) []*websocket.Conn {
	if len(indices) == 0 {
		return conns
	}

	size := len(conns) - len(indices)
	ret := make([]*websocket.Conn, size)
	insert := 0

	for i := 0; i < len(conns); i++ {
		if _, ok := indices[i]; ok {
			continue
		}

		ret[insert] = conns[i]
		insert = insert + 1
	}

	return ret
}

func (app *StatusApp) RemoveDevice(ipAddress string) error {
	app.Lock()
	defer app.Unlock()

	if _, ok := app.Devices[ipAddress]; !ok {
		return fmt.Errorf("Device %s does not exist", ipAddress)
	}

	delete(app.Devices, ipAddress)

	return app.Save()
}

func (app *StatusApp) AddDevice(ipAddress string, name string) error {
	app.Lock()
	defer app.Unlock()

	// Check if device exists
	if _, ok := app.Devices[ipAddress]; ok {
		return fmt.Errorf("Device with that IP already listed")
	}

	app.Devices[ipAddress] = DeviceInfo{name}

	return app.Save()
}

// Sometimes it's just too easy
