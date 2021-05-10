package subprocess

import (
	"fmt"
	"os"
)

type ColoredWriter struct {
	out *os.File
}

func (this ColoredWriter) Write(b []byte) (int, error) {
	fmt.Fprintf(this.out, "\033[0;34m%s\033[0m", b)
	return len(b), nil
}
