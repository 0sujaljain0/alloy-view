package querier

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/0sujaljain0/alloy-view/pkg/state"
	"github.com/0sujaljain0/alloy-view/pkg/utils"
)

type AlloyApiQuerier interface {
	Init(state.AppState) error
	buildInventory(nodes []state.AlloyNode) error
	SearchTarget(string) error
}

type BaseAlloyApiQuerier struct {
	nodes      []state.AlloyNode
	logger     *slog.Logger
	Components *AlloyComponents
}

type AlloyComponents struct {
	mutex      *sync.Mutex
	components map[uint32]AlloyComponent
}

func (c *AlloyComponents) Len() int {
	return len(c.components)
}

func (c *AlloyComponents) AddComponent(comp AlloyComponent, logger *slog.Logger) {
	if comp == nil {
		return
	}

	idKey := utils.HashString(comp.GetName())
	_, found := c.components[idKey]

	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !found {
		logger.Info(fmt.Sprintf("discovered: %s", comp.GetName()))
		c.components[idKey] = comp
		return
	} 

	logger.Debug(fmt.Sprintf("already discovered: %s, so skipping", comp.GetName()))
}

func (c *AlloyComponents) GetComponents() []AlloyComponent{ 
	var values []AlloyComponent = make([]AlloyComponent, 0)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, value := range c.components {
		values = append(values, value)
	}

	return values
}

func (q *BaseAlloyApiQuerier) buildComponentInventory() error {
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	for _, node := range q.nodes {
		wg.Add(1)
		go addNodeComponents(node, q.Components, wg, q.logger)
	}
	wg.Wait()
	return nil
}

func addNodeComponents(node state.AlloyNode, comps *AlloyComponents, wg *sync.WaitGroup, logger *slog.Logger) error {
	defer wg.Done()

	resp, err := http.Get(fmt.Sprintf("http://%s/api/v0/web/components", node.GetEndpoint()))
	if err != nil {
		resp.Body.Close()
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data []*BaseAlloyComponent
	if err = json.Unmarshal(body, &data); err != nil {
		return err
	}
	for _, ele := range data {
		// logger.Info(fmt.Sprintf("name of the component: %s", ele.Name))
		comps.AddComponent(ele, logger)
	}

	logger.Info(fmt.Sprintf("len of components: %d", comps.Len()))


	return nil
}

type HealthState struct {
	State       string `json:"state"`
	Message     string `json:"message"`
	UpdatedTime string `json:"updatedTime"`
}

func (s *HealthState) String() string {
	return fmt.Sprintf("{ %s | %s | %s }", s.State, s.Message, s.UpdatedTime)
}

type AlloyComponent interface {
	GetHealth() *HealthState
	GetName() string
}

type BaseAlloyComponent struct {
	Name   string       `json:"localID"`
	Type   string       `json:"name"`
	Health *HealthState `json:"health"`
}

func (b *BaseAlloyComponent) GetHealth() *HealthState { return b.Health }
func (b *BaseAlloyComponent) GetName() string { return b.Name }

//TODO : THIS NEEDS TO BE IMPLEMENTED FIRST
func (b *BaseAlloyComponent) concretizeAlloyComponent() AlloyComponent {
	return nil
}

type PrometheusScrapeAlloyComponent struct {
	*BaseAlloyComponent
	Scrapes Scrapes
}


//TODO : THIS NEEDS TO BE IMPLEMENTED FIRST
type Scrape struct {
}

type Scrapes []*Scrape
