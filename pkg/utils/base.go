package utils

import (
	"fmt"
	"hash/fnv"
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

func HashString(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}
