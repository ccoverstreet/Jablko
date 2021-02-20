package types

import (
	"net/http"
)

type JablkoInterface interface {
	SyncConfig(string) error
	SendMessage(string) error
	GetFlagValue(string) bool
}

type JablkoMod interface{
	ConfigStr() ([]byte, error)
	ModuleCardConfig() string
	WebHandler(http.ResponseWriter, *http.Request)
	UpdateConfig([]byte) error
	SourcePath() string
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
