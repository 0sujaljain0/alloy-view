package handler

import (
	"context"
	"net/http"

	"github.com/0sujaljain0/alloy-view/pkg/view"
	"github.com/0sujaljain0/alloy-view/pkg/view/components"
)

func (h *HandlerClustered) ServeHomePage(res http.ResponseWriter, req *http.Request) {
	view.Home(h.State).Render(context.Background(), res)
}

func (h *HandlerClustered) ClusterInfoComp(res http.ResponseWriter, req *http.Request) {
	components.ClusterInfo(h.State.GetNodes()).Render(context.Background(), res)
}
