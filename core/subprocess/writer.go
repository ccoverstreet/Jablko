package subprocess

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type ColoredWriter struct {
	prefix string
	out    *os.File
}

func (this ColoredWriter) Write(b []byte) (int, error) {
	fmt.Fprintf(this.out, "\033[0;34m%s: %s\033[0m", this.prefix, b)
	return len(b), nil
}

type SubprocessWriter struct {
	sync.Mutex
	JMODName string
	fileName string
	logFile  os.File

	curLength int
}

func CreateSubprocessWriter(JMODName string) *SubprocessWriter {
	writer := new(SubprocessWriter)
	writer.JMODName = JMODName
	writer.fileName = strings.ReplaceAll(JMODName, "/", "_")

	return writer
}
