package main

import (
	"os"
	"fmt"

	"github.com/0sujaljain0/alloy-view/pkg/config"
	"github.com/0sujaljain0/alloy-view/pkg/utils"
	"github.com/0sujaljain0/alloy-view/pkg/web"
)

func main() {
	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	logger, closer := utils.NewLogger("logs.log")
	defer closer()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Error("config.yaml not found")
		os.Exit(1)
	}

	logger.Info("Reading config.yaml")
	cfg, err := config.ParseConfig(data)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("%s", *cfg))

	server, err := web.ConfigureServer(8093, "dev.testing.server", cfg, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("%s", server))
	err = server.Start()
	if err != nil {
		logger.Error(err.Error())
	}
}
