package state

import (
	"fmt"
	"log/slog"

	"github.com/0sujaljain0/alloy-view/pkg/config"
	"github.com/0sujaljain0/alloy-view/pkg/utils"
)

type AppState interface {
	String() string
	populateNodes() error
	GetNodes() []AlloyNode
	InitState() error
}

type BaseAppState struct {
	mode   string
	Nodes  []AlloyNode
	SD     config.ServiceDiscovery
	logger *slog.Logger
}

type ClusteredAppState struct {
	BaseAppState
}

func NewClusteredAppState(sd config.ServiceDiscovery, mode string, logger *slog.Logger) *ClusteredAppState {
	return &ClusteredAppState{
		BaseAppState: BaseAppState{
			Nodes:  make([]AlloyNode, 0),
			SD:     sd,
			mode:   mode,
			logger: logger,
		},
	}
}

func (s *ClusteredAppState) String() string {
	return "{ not implemented the String() function for ClusteredAppState correctly }"
}
func (s *ClusteredAppState) GetNodes() []AlloyNode { return s.Nodes }
func (s *ClusteredAppState) populateNodes() error {
	switch mode := s.SD.GetMode(); mode {
	case "k8s":
		sdConf, ok := s.SD.(*config.K8sServiceDiscovery)
		if !ok {
			return fmt.Errorf("mode:cluster, but s.SD not correctly typed")
		}

		eps := utils.GetEPSFromK8sSvc(*sdConf.Namespace, *sdConf.Service, s.logger)
		for _, ep := range eps {
			s.Nodes = append(s.Nodes, &K8sAlloyNode{
				PodName:  ep.PodName,
				NodeName: ep.NodeName,
				IP:       ep.IP,
			})
		}
		return nil

	default:
		return fmt.Errorf("Invalid Service Discovery Mode: %s", mode)
	}
}
func (s *ClusteredAppState) InitState() error {
	err := s.populateNodes()
	if err != nil {
		return err
	}

	return nil
}
