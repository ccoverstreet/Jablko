package types

import (
	"net/http"
)

type JablkoMod interface{
	Card(*http.Request) string
	WebHandler(http.ResponseWriter, *http.Request)
}

func StructToMod(inputStruct JablkoMod) JablkoMod {
	return inputStruct
}
