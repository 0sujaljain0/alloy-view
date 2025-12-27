package querier

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/0sujaljain0/alloy-view/pkg/state"
)

type ClusterAlloyApiQuerier struct {
	BaseAlloyApiQuerier
}

func NewClusterAlloyApiQuerier(logger *slog.Logger) AlloyApiQuerier {
	logger.Debug("ClusterAlloyApiQuerier created, returning")
	return &ClusterAlloyApiQuerier{
		BaseAlloyApiQuerier: BaseAlloyApiQuerier{
			logger: logger,
			nodes:  make([]state.AlloyNode, 0),
			Components: &AlloyComponents{
				mutex:      &sync.Mutex{},
				components: make([]AlloyComponent, 0),
			},
		},
	}
}

func (q * ClusterAlloyApiQuerier) Init(s state.AppState) error {
	if err := q.buildInventory(s.GetNodes()); err != nil {
		return err
	}
	return nil
}

func (q *ClusterAlloyApiQuerier) buildInventory(nodes []state.AlloyNode) error {
	if len(nodes) == 0 {
		return fmt.Errorf("no nodes supplied")
	}
	q.nodes = nodes

	if err := q.buildComponentInventory(); err != nil {
		return fmt.Errorf("while building component inventory")
	}

	return nil
}

func (q *ClusterAlloyApiQuerier) SearchTarget(keyword string) error {
	return fmt.Errorf("not implemented yet")
}
