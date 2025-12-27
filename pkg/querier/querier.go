package querier

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/0sujaljain0/alloy-view/pkg/state"
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
	components []AlloyComponent
}

func (c *AlloyComponents) AddComponent(comp AlloyComponent) {
	if comp == nil {
		return
	}
	c.mutex.Lock()
	c.components = append(c.components, comp)
	c.mutex.Unlock()
}

func (c *AlloyComponents) GetComponents() []AlloyComponent { return c.components }

func (q *BaseAlloyApiQuerier) buildComponentInventory() error {
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	for _, node := range q.nodes[:1] {
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

	// FIX: THIS NEEDS TO BE FIXED, UNMARSHALLING LOGIC IS NOT SEGREGATING THE OBJECTS INSIDE THE COMPONENTS ARRAY IN
	// THE RESPONSE.
	var data []map[string]any
	err = json.Unmarshal(body, &data)
	logger.Info(fmt.Sprintf("%d", len(data[0])))
	logger.Info(fmt.Sprintf(" len -> %d", len(data)))
	for _, ele:= range data {
		for key, value := range ele {
			logger.Info(fmt.Sprintf("%s -> %+v", key, value))
		}
	}
	comps.AddComponent(nil)
	if err != nil {
		return err
	}

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
}

type BaseAlloyComponent struct {
	Name   string       `json:"localID"`
	Type   string       `json:"name"`
	Health *HealthState `json:"health"`
}

func (b *BaseAlloyComponent) GetHealth() *HealthState { return b.Health }

type PrometheusScrapeAlloyComponent struct {
	BaseAlloyComponent
	Scrapes Scrapes
}

type Scrape struct {
}

type Scrapes []*Scrape
