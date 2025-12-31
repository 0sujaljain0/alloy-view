package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"github.com/0sujaljain0/alloy-view/pkg/view"
	"github.com/0sujaljain0/alloy-view/pkg/view/components"
)

// /////////////// PUBLIC SERVING ///////////////////////////
func (h *HandlerClustered) ServeHomePage(res http.ResponseWriter, req *http.Request) {
	view.Home(h.State).Render(context.Background(), res)
}

func (h *HandlerClustered) ServeComponentsPage(res http.ResponseWriter, req *http.Request) {
	view.ComponentsPage(h.apiQuerier.GetComponents()).Render(context.Background(), res)
}
func (h *HandlerClustered) ServeNodesInfoPage(res http.ResponseWriter, req *http.Request) {
	h.State.RefreshState()
	view.NodeInfoPage(h.State.GetNodes()).Render(context.Background(), res)
}

func (h *HandlerClustered) ServeSearchTargetPage(res http.ResponseWriter, req *http.Request) {
	view.TargetSearchPage().Render(context.Background(), res)
}

func (h *HandlerClustered) TargetSearchSubmitHandler(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		h.logger.Error("error while parsing the form")
	}

	components.HtmlDataWrapper(fmt.Sprintf("%s, %s", req.FormValue("typeOfSearch"), req.FormValue("search_keyword"))).Render(context.Background(), res)
}

////////////////////////////////////////////////////////////

//go:generate stringer -type=NodeHealthStatus
type NodeHealthStatus uint32

const (
	HEALTHY NodeHealthStatus = iota
	UNHEALTHY
	NOT_REACHABLE
)
/////////////////////////////////////////////////////////////
// /////////////////// INTERNAL SERVING /////////////////////
/////////////////////////////////////////////////////////////
func (h *HandlerClustered) ClusterInfoComp(res http.ResponseWriter, req *http.Request) {
	h.State.RefreshState()
	components.ClusterInfo(h.State.GetNodes()).Render(context.Background(), res)
}

func (h *HandlerClustered) ServeNodeHealthIndicator(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	endpoint := query.Get("alloy_endpoint")
	size := query.Get("size") // "small" or "big"

	if endpoint == "" {
		http.Error(res, "alloy_endpoint is required", http.StatusBadRequest)
		return
	}

	// Check health status
	healthy := h.checkNodeHealth(endpoint)

	// Render appropriate indicator
	if size == "big" {
		components.HealthIndicatorBig(healthy).Render(context.Background(), res)
	} else {
		// Default to small
		components.HealthIndicatorSmall(healthy).Render(context.Background(), res)
	}
}

func (h *HandlerClustered) checkNodeHealth(endpoint string) bool {
	resp, err := http.Get("http://" + endpoint + "/-/healthy")
	if err != nil {
		h.logger.Error("health check failed", "endpoint", endpoint, "error", err.Error())
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("failed to read health response", "error", err.Error())
		return false
	}

	h.logger.Debug("health check response", "endpoint", endpoint, "status", resp.StatusCode, "body", string(body))

	return resp.StatusCode == 200
}

func (h *HandlerClustered) ServeComponentsFilter(res http.ResponseWriter, req *http.Request) {
    // 1. Parse query parameters
    search := req.URL.Query().Get("search")
    searchMode := req.URL.Query().Get("search_mode") 
    typeFilter := req.URL.Query().Get("type_filter")
    
    // 2. Get all components
    allComponents := h.apiQuerier.GetComponents().GetComponents()
    
    // 3. Apply filters
    filteredComponents := filterComponents(allComponents, search, searchMode, typeFilter)
    
    // 4. Render the filtered list
    view.ComponentsList(filteredComponents).Render(context.Background(), res)
}

////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
