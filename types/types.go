package types

import (
	"net/http"
)

type JablkoInterface interface {
	Tester()
	SyncConfig(string)
}

type JablkoMod interface{
	ConfigStr() ([]byte, error)
	Card(*http.Request) string
	WebHandler(http.ResponseWriter, *http.Request)
}

func StructToMod(inputStruct JablkoMod) JablkoMod {
	return inputStruct
}
