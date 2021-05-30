package subprocess

import (
	"fmt"
	"os"
)

type ColoredWriter struct {
	prefix string
	out    *os.File
}

func (this ColoredWriter) Write(b []byte) (int, error) {
	fmt.Fprintf(this.out, "\033[0;34m%s: %s\033[0m", this.prefix, b)
	return len(b), nil
}
