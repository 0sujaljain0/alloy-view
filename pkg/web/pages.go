package web

import (
	"context"
	"net/http"

	"github.com/0sujaljain0/alloy-view/pkg/view"
)

func (s *Server) ServeHomePage(res http.ResponseWriter, req *http.Request) {
	view.Home().Render(context.Background(), res)
}
