package config

import (
	"fmt"
)

type ServiceDiscovery interface {
	GetMode() string
	Validate() error
}

type BaseServiceDiscovery struct {
	Mode string `yaml:"sdMode"`
}

func (sd *BaseServiceDiscovery) String() string { return fmt.Sprintf("[ serviceDiscoveryMode: %s ]", sd.Mode) }

type K8sServiceDiscovery struct {
	BaseServiceDiscovery `yaml:",inline"`
	Service *string `yaml:"serviceName"`
	Namespace *string `yaml:"namespace"`
}

func (sd *K8sServiceDiscovery) String() string {
	return fmt.Sprintf("{ %s - [%s/%s] }", &sd.BaseServiceDiscovery, *sd.Namespace, *sd.Service)
}
func (sd *K8sServiceDiscovery) GetMode() string { return sd.Mode }
func (sd *K8sServiceDiscovery) Validate() error {
	if sd.Service == nil {
		return fmt.Errorf("invalid config: must provide serviceDiscovery.serviceName")
	}
	if sd.Namespace == nil {
		return fmt.Errorf("invalid config: must provide serviceDiscovery.namespace")
	}

	return nil
}
