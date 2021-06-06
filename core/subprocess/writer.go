package subprocess

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type SubprocessWriter struct {
	sync.Mutex
	JMODName string
	fileName string
	logFile  *os.File
	curDay   int
}

func CreateSubprocessWriter(JMODName string) (*SubprocessWriter, error) {
	writer := new(SubprocessWriter)
	writer.JMODName = JMODName
	writer.fileName = fmt.Sprintf("./log/%s_%s.log", strings.ReplaceAll(writer.JMODName, "/", "_"), time.Now().Format("2006-01-02"))

	logFile, err := os.OpenFile(writer.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	writer.logFile = logFile

	writer.curDay = time.Now().Day()

	return writer, nil
}

func (writer *SubprocessWriter) Write(b []byte) (int, error) {
	writer.Lock()
	defer writer.Unlock()
	if writer.curDay != time.Now().Day() {
		log.Info().
			Str("jmodName", writer.JMODName).
			Msg("Cycling log file for JMOD")

		writer.logFile.Close() // Closing old file

		writer.curDay = time.Now().Day()
		writer.fileName = fmt.Sprintf("./log/%s_%s.log", strings.ReplaceAll(writer.JMODName, "/", "_"), time.Now().Format("2006-01-02"))

		newLogFile, err := os.OpenFile(writer.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}

		writer.logFile = newLogFile
	}

	fmt.Printf("\033[0;34m%s\033[0m", b)
	writer.logFile.Write(b)

	return len(b), nil
}
