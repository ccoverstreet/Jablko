// Jablko Logging
// Cale Overstreet
// Apr. 13, 2021

/* 
Structured logging for Jablko. The idea for making my
own structured logging is to learn the concepts and
tailor the logging to my specific needs
*/

import (
	"fmt"
	"os"
)

// Global
type JLogger {
	Output os.File
}

var logger = JLogger{os.Stdout}

func Info(v interface{}...) {
	
}
