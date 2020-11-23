package types

import (
	"net/http"
)

type JablkoInterface interface {
	Tester()
	SyncConfig(string)
	SendMessage(string) error
}

type JablkoMod interface{
	ConfigStr() ([]byte, error)
	Card(*http.Request) string
	WebHandler(http.ResponseWriter, *http.Request)
}

func StructToMod(inputStruct JablkoMod) JablkoMod {
	return inputStruct
}

type UserData struct {
	Id int
	Username string
	Password string
	FirstName string
	Permissions int
}

type SessionHolder struct {
	Id int
	Cookie string
	Username string
	FirstName string
	Permissions int
	CreationTime int
}
