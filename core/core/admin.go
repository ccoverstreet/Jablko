package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func httpErrorHandler(w http.ResponseWriter, aErr *HTTPError) {
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
}

func AdminFuncHandler(w http.ResponseWriter, r *http.Request, core *JablkoCore) {
	fun := mux.Vars(r)["func"]

	handler, ok := ADMINFUNCMAP[fun]

	if !ok {
		httpErrorHandler(w, CreateHTTPError(400,
			fmt.Sprintf("Invalid admin func '%s' requested", fun), nil))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpErrorHandler(w,
			CreateHTTPError(400, "Unable to read request body", err))
	}

	res, herr := handler(body, core)
	if herr != nil {
		httpErrorHandler(w, herr)
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
		httpErrorHandler(w,
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

	return struct{}{}, nil
}
