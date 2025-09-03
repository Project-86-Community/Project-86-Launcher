package p86l

import (
	"os"

	"github.com/rs/zerolog"
)

type LogModel struct {
	logger  zerolog.Logger
	logFile *os.File
}

func (l *LogModel) Logger() zerolog.Logger {
	return l.logger
}

func (l *LogModel) LogFile() *os.File {
	return l.logFile
}

func (l *LogModel) SetLogger(logger zerolog.Logger) {
	l.logger = logger
}

func (l *LogModel) SetLogFile(logFile *os.File) {
	l.logFile = logFile
}

// -- common --

func (l *LogModel) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
