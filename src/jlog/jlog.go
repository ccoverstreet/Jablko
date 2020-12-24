// jlog.go: Jablko Logging Modules
// Cale Overstreet
// 2020/12/24
// This module contains the logging functions for Jablko. The
// primary functions are jlog.Printf, jlog.Warnf, jlog.Errorf.
// These functions output text in a color coded manner. This
// feature was added to make debugging and development easier.


package jlog

import (
	"fmt"
	"time"
)

var colorMapANSI = map[string]string {
	"red": "\033[0;31m",
	"yellow": "\033[0;33m",
	"reset": "\033[0m",
}

func prefix() string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "]: "
}

func Printf(format string, args ...interface{}) {
	fmt.Printf(prefix() + format, args...)
}

func Warnf(format string, args ...interface{}) {
	fmt.Printf(colorMapANSI["yellow"])
	fmt.Printf(prefix() + format, args...)
	fmt.Printf(colorMapANSI["reset"])
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf(colorMapANSI["red"])
	fmt.Printf(prefix() + format, args...)
	fmt.Printf(colorMapANSI["reset"])
}
