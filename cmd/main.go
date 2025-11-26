package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"

	"gopkg.in/yaml.v3"
)

type ParticipantSDMode uint8

const (
	IpSD ParticipantSDMode = iota
	k8sSD
)

func (m ParticipantSDMode) String() string {
	switch m {
	case IpSD:
		return "IpSD"
	case k8sSD:
		return "k8SD"
	default:
		return "Participation type invalid"
	}
}

type AlloyModeConfig interface {
	ConfigInfo() string
}

type ParticipantSD struct {
	participantMode ParticipantSDMode
	sdEndpoint      string
}

type AlloyClusterModeConfig struct {
	SD *ParticipantSD `yaml:"participantSD"`
}

func (c *AlloyClusterModeConfig) ConfigInfo() string {
	return fmt.Sprintf("{ Mode: %s | sdEndpoint: %s }", c.SD.participantMode, c.SD.sdEndpoint)
}

type AlloyIndividualModeConfig struct {
	SD *ParticipantSD `yaml:"participantSD"`
}

func (c *AlloyIndividualModeConfig) ConfigInfo() string {
	return fmt.Sprintf("{ Mode: %s | sdEndpoint: %s }", c.SD.participantMode, c.SD.sdEndpoint)
}

type RawConfigWrapper struct {
	Mode             string                     `yaml:"mode"`
	ClusterConfig    *AlloyClusterModeConfig    `yaml:"cluster,omitempty"`
	IndividualConfig *AlloyIndividualModeConfig `yaml:"individual,omitempty"`
}

func ParseConfig(configData []byte, logger *slog.Logger) (AlloyModeConfig, error) {
	config := &RawConfigWrapper{}
	err := yaml.Unmarshal(configData, config)
	if err != nil {
		logger.Error(err.Error())
	}
	switch config.Mode {
	case "cluster":
		return config.ClusterConfig, nil
	case "individual":
		return config.IndividualConfig, nil
	default:
		return nil, fmt.Errorf("Mode not available: %s", config.Mode)
	}
}

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	logFile, err := os.OpenFile("logs.logs", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	defer logFile.Close()

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
	config, err := ParseConfig(data, logger)
	if err != nil {
		panic(err)
	}
	logger.Info(config.ConfigInfo())
}
