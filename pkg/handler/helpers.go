package handler

import (
	"strings"
	"strconv"
	"github.com/0sujaljain0/alloy-view/pkg/querier"
)

///////////////////////////////////////////////////////////////////
// this section is for the Handler.ServeComponentsFilter function
///////////////////////////////////////////////////////////////////
func filterComponents(components []querier.AlloyComponent, search, searchMode, typeFilter string) []querier.AlloyComponent {
    var result []querier.AlloyComponent
    
    for _, comp := range components {
        // Type filter
        if typeFilter != "all" && comp.GetType() != typeFilter {
            continue
        }
        
        // Search filter
        if search != "" && !matchesSearch(comp, search, searchMode) {
            continue
        }
        
        result = append(result, comp)
    }
    
    return result
}

func matchesSearch(comp querier.AlloyComponent, search, mode string) bool {
    search = strings.ToLower(search)
    
    switch mode {
    case "component":
        return strings.Contains(strings.ToLower(comp.GetName()), search)
    case "target":
        // Only search if it's a Prometheus component
        if promComp, ok := comp.(*querier.PrometheusScrapeAlloyComponent); ok {
            targetStr := strconv.Itoa(len(promComp.Scrapes))
            return targetStr == search
        }
        return false
    default:
        return true
    }
}
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
