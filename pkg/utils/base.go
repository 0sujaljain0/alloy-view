package utils

import (
	"fmt"
	"log/slog"
	"os"
)

func NewLogger(logFilePath string) (*slog.Logger, func()) {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	closeFunc := func() {
		err := logFile.Close()
		if err != nil {
			panic(fmt.Errorf("while closing the file: %+v", err))
		}
	}

	if err != nil {
		panic(fmt.Sprintf("log file not initialized: %s", err))
	}

	return slog.New(slog.NewTextHandler(logFile, nil)), closeFunc
}
