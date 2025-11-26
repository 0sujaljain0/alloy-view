package config

import (
	"fmt"
	"log/slog"
	"gopkg.in/yaml.v3"
)

type ParticipantSDMode string

type ParticipantSD struct {
	ParticipantMode ParticipantSDMode `yaml:"sdMode"`
	SDEndpoint      string            `yaml:"endpoint"`
}

type AlloyModeConfig interface {
	ConfigInfo() string
}

type AlloyClusterModeConfig struct {
	SD *ParticipantSD `yaml:"participantSD"`
}

func (c *AlloyClusterModeConfig) ConfigInfo() string {
	return fmt.Sprintf("{ Mode: %s | sdEndpoint: %s }", c.SD.ParticipantMode, c.SD.SDEndpoint)
}

type AlloyIndividualModeConfig struct {
	SD *ParticipantSD `yaml:"participantSD"`
}

func (c *AlloyIndividualModeConfig) ConfigInfo() string {
	return fmt.Sprintf("{ Mode: %s | sdEndpoint: %s }", c.SD.ParticipantMode, c.SD.SDEndpoint)
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
