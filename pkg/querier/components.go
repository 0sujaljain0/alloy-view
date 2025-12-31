package querier

import (
	"fmt"
	"regexp"

	"github.com/0sujaljain0/alloy-view/pkg/state"
)

type AlloyComponent interface {
	GetHealth() *HealthState
	GetName() string
	GetType() string
	getBaseComponent() *BaseAlloyComponent
	PopulateStruct([]state.AlloyNode) error
}

type BaseAlloyComponent struct {
	Name           string       `json:"localID"`
	Type           string       `json:"name"`
	Health         *HealthState `json:"health"`
}

func (b *BaseAlloyComponent) GetType() string         					   { return b.Type }
func (b *BaseAlloyComponent) GetHealth() *HealthState 					   { return b.Health }
func (b *BaseAlloyComponent) GetName() string         					   { return b.Name }
func (b *BaseAlloyComponent) PopulateStruct(nodes []state.AlloyNode) error { return nil}
func (b *BaseAlloyComponent) getBaseComponent() *BaseAlloyComponent { return b }

var (
	rePrometheusScrape = regexp.MustCompile(".*prometheus.scrape.*")
)

// TODO: 2. Create a Concrete Component for discovery.relabel also
type PrometheusScrapeAlloyComponent struct {
	*BaseAlloyComponent
	Scrapes Scrapes
}

// TODO: 1. Figure out a way collecting all the scrape targets from the all the nodes for a "prometheus.scrape component".
func (ps *PrometheusScrapeAlloyComponent) PopulateStruct(node []state.AlloyNode) error { 
	return nil 
}

type Scrape struct {
}

type Scrapes []*Scrape


type HealthState struct {
	State       string `json:"state"`
	Message     string `json:"message"`
	UpdatedTime string `json:"updatedTime"`
}

func (s *HealthState) String() string {
	return fmt.Sprintf("{ %s | %s | %s }", s.State, s.Message, s.UpdatedTime)
}

func concretizeAlloyComponent(comp AlloyComponent) AlloyComponent {
	switch {
	case rePrometheusScrape.MatchString(comp.GetType()):
		component := &PrometheusScrapeAlloyComponent{
			BaseAlloyComponent: comp.getBaseComponent(),
			Scrapes:            make(Scrapes, 0),
		}
		return component
	default:
		return comp
	}
}
