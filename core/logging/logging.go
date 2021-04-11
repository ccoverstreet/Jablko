// Jablko Core Logging
// Cale Overstreet
// Mar. 30, 2021

// Contains custom io logger 

package logging

import (
	"fmt"
	"time"
)

type JablkoLogger struct{}

func (writer JablkoLogger) Write(bytes []byte) (int, error) {
	return fmt.Printf("[" + time.Now().Format(time.RFC1123) + "]: " + string(bytes))
}
