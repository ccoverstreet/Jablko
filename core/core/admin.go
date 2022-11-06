package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ccoverstreet/Jablko/core/process"
	"github.com/gorilla/mux"
)

func HTTPErrorHandler(w http.ResponseWriter, aErr *HTTPError) {
	b, err := json.Marshal(aErr)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(aErr.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

var ADMINFUNCMAP = map[string]func(body []byte, core *JablkoCore) (interface{}, *HTTPError){
	"asd":        exampleAdminFunc,
	"getModList": getModListHandler,
	"removeMod":  removeModHandler,
	"addMod":     addModHandler,
	"updateMod":  updateModHandler,
}

func AdminFuncHandler(w http.ResponseWriter, r *http.Request, core *JablkoCore) {
	fun := mux.Vars(r)["func"]

	handler, ok := ADMINFUNCMAP[fun]

	if !ok {
		HTTPErrorHandler(w, CreateHTTPError(400,
			fmt.Sprintf("Invalid admin func '%s' requested", fun), nil))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		HTTPErrorHandler(w,
			CreateHTTPError(400, "Unable to read request body", err))
		return
	}

	res, herr := handler(body, core)
	if herr != nil {
		HTTPErrorHandler(w, herr)
		return
	}

	JSONResponse(w, res)
}

type HTTPError struct {
	StatusCode int    `json:"-"`
	Err        string `json:"err"`
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("<%d> %v", e.StatusCode, e.Err)
}

func CreateHTTPError(statusCode int, message string, err error) *HTTPError {
	return &HTTPError{statusCode, fmt.Sprintf("%v: %v", message, err)}
}

func JSONResponse(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		HTTPErrorHandler(w,
			CreateHTTPError(500, "Unable to marshal handler JSON res", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func exampleAdminFunc(b []byte, core *JablkoCore) (interface{}, *HTTPError) {
	return nil, &HTTPError{400,
		fmt.Sprintf("Some error the client is likely responsible for: %v")}
}

func getModListHandler(b []byte, core *JablkoCore) (interface{}, *HTTPError) {
	return &core.PMan, nil
}

// Adds the mod to the manager, runs its update function, and tries to start it
func addModHandler(b []byte, core *JablkoCore) (interface{}, *HTTPError) {
	procName := process.DetermineProcName(b)
	procTag := process.DetermineProcTag(b)

	fmt.Println(procName, procTag)

	err := core.PMan.AddMod(procName, b)
	if err != nil {
		return struct{}{},
			CreateHTTPError(500, "Unable to add mod", err)
	}

	err = core.PMan.UpdateMod(procName, procTag)
	if err != nil {
		return struct{}{},
			CreateHTTPError(500, "Unable to run update function for mod", err)
	}

	err = core.SaveConfig()
	if err != nil {
		return struct{}{},
			CreateHTTPError(500, "Unable to save config", err)
	}

	err = core.PMan.StartMod(procName)
	if err != nil {
		return struct{}{},
			CreateHTTPError(500, "Unable to save config", err)
	}

	return struct{}{}, nil
}

func removeModHandler(b []byte, core *JablkoCore) (interface{}, *HTTPError) {
	input := struct {
		Name string `json:"name"`
	}{}

	err := json.Unmarshal(b, &input)
	if err != nil {
		return nil, CreateHTTPError(500, "Unable to marshal JSON for removeMod", err)
	}

	err = core.PMan.RemoveMod(input.Name)
	if err != nil {
		return nil,
			CreateHTTPError(500, "Unable to remove mod", err)
	}

	err = core.SaveConfig()
	if err != nil {
		return struct{}{},
			CreateHTTPError(500, "Unable to save config", err)
	}

	return struct{}{}, nil
}

func updateModHandler(b []byte, core *JablkoCore) (interface{}, *HTTPError) {
	input := struct {
		Name string `json:"name"`
		Tag  string `json:"tag"`
	}{}

	err := json.Unmarshal(b, &input)

	if err != nil {
		return struct{}{},
			CreateHTTPError(400, "Invalid input for updateModHandler", err)
	}

	err = core.PMan.UpdateMod(input.Name, input.Tag)
	if err != nil {
		return struct{}{},
			CreateHTTPError(500, "Unable to update mod", err)
	}

	err = core.PMan.StopMod(input.Name)
	if err != nil {
		return struct{}{},
			CreateHTTPError(500, "Unable to stop mod", err)
	}

	err = core.PMan.StartMod(input.Name)
	if err != nil {
		return struct{}{},
			CreateHTTPError(500, "Unable to start mod", err)
	}

	return struct{}{}, nil
}
