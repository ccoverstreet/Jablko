// Jablko Database Handler
// Cale Overstreet
// May 6, 2021

/*
All accesses to Jablko's database should go through
this handler. Handles creation and state management 
for user logins, registered pmods.
*/

package database

import (
	"database/sql"

	"github.com/rs/zerolog/log"
)

type DatabaseHandler struct {
	i
}
