package process

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type JMODWriter struct {
	sync.Mutex
	ImageName        string
	ImageNameCleaned string
	logFile          *os.File
	curDay           int
}

func OpenLogFileFromName(imageName string) (*os.File, error) {
	logFileName := fmt.Sprintf("log/%s_%s.log", CleanImageName(imageName), time.Now().Format("2006-01-02"))
	return os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func CreateJMODWriter(imageName string) (*JMODWriter, error) {
	logFile, err := OpenLogFileFromName(imageName)
	if err != nil {
		return nil, err
	}

	writer := &JMODWriter{
		sync.Mutex{},
		imageName,
		CleanImageName(imageName),
		logFile,
		time.Now().Day(),
	}

	return writer, nil
}

func (writer *JMODWriter) Write(b []byte) (int, error) {
	writer.Lock()
	defer writer.Unlock()

	if writer.curDay != time.Now().Day() {
		log.Info().
			Str("imageName", writer.ImageName).
			Msg("Cyling JMOD log file")

		writer.logFile.Close()

		newFile, err := OpenLogFileFromName(writer.ImageName)
		if err != nil {
			log.Error().
				Str("imageName", writer.ImageName).
				Msg("Failed to rotate log file")

			return 0, err
		}

		writer.logFile = newFile
		writer.curDay = time.Now().Day()
	}

	fmt.Printf("\033[0;34m%s: %s\033[0m", writer.ImageName, b)
	writer.logFile.Write(b)

	return len(b), nil
}
