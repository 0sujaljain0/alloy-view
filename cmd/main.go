package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"

	"github.com/0sujaljain0/alloy-view/pkg/config"
)


func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	logFile, err := os.OpenFile("logs.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	defer func() {
		err := logFile.Close()
		if err != nil {
			panic(fmt.Errorf("while closing the file: %+v", err))
		}
	}()

	if err != nil {
		panic(fmt.Sprintf("log file not initialized: %s", err))
	}

	logger := slog.New(slog.NewTextHandler(logFile, nil))

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Error("config.yaml not found")
		os.Exit(1)
	}

	logger.Info("Reading config.yaml")
	config, err := config.ParseConfig(data, logger)
	if err != nil {
		panic(err)
	}
	logger.Info(config.ConfigInfo())
}
