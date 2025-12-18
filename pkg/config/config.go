package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type AlloyConfig interface {
	GetMode() string
	Validate() error
	String() string
}

type BaseConfig struct {
	Mode string `yaml:"mode"`
}

func (c *BaseConfig) String() string {
	return fmt.Sprintf("[mode: %s]", c.Mode)
}

type ClusterConfig struct {
	BaseConfig `yaml:",inline"`
	SD         ServiceDiscovery
}

func (c *ClusterConfig) String() string {
	return fmt.Sprintf("{ %s - %s }", &c.BaseConfig, c.SD)
}
func (c *ClusterConfig) GetMode() string { return c.Mode }
func (c *ClusterConfig) Validate() error {
	if err := c.SD.Validate(); err != nil {
		return err
	}
	return nil
}

func (c *ClusterConfig) UnmarshalYAML(value *yaml.Node) error {
	type plain ClusterConfig

	type rawWrapper struct {
		*plain `yaml:",inline"`
		RawSD  yaml.Node `yaml:"serviceDiscovery"`
		// add more fields.
	}
	wrapper := rawWrapper{plain: (*plain)(c)}
	if err := value.Decode(&wrapper); err != nil {
		return err
	}

	{ // ServiceDiscovery
		type sdDisc struct {
			SDMode string `yaml:"sdMode"`
		}
		var disc sdDisc

		if err := wrapper.RawSD.Decode(&disc); err != nil {
			return err
		}
		switch disc.SDMode {
		case "k8s":
			c.SD = &K8sServiceDiscovery{}
		default:
			return fmt.Errorf("unknown sdMode: %s", disc.SDMode)
		}

		if err := wrapper.RawSD.Decode(c.SD); err != nil {
			return err
		}
	}
	return nil
}

func ParseConfig(data []byte) (*AlloyConfig, error) {
	var node yaml.Node

	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, err
	}

	var baseConfig BaseConfig

	if err := node.Decode(&baseConfig); err != nil {
		return nil, err
	}

	var finalConfig AlloyConfig

	switch baseConfig.Mode {
	case "cluster":
		var cfg ClusterConfig
		if err := node.Decode(&cfg); err != nil {
			return nil, err
		}
		finalConfig = &cfg
	default:
		return nil, fmt.Errorf("Invalid Mode: %s", baseConfig.Mode)
	}

	if err := finalConfig.Validate(); err != nil {
		return nil, err
	}

	return &finalConfig, nil
}
