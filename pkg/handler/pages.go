package handler

import (
	"context"
	"net/http"

	"github.com/0sujaljain0/alloy-view/pkg/view"
)

func (h *HandlerClustered) ServeHomePage(res http.ResponseWriter, req *http.Request) {
	view.Home(h.State).Render(context.Background(), res)
}

