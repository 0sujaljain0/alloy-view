package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/0sujaljain0/alloy-view/pkg/config"
	"github.com/0sujaljain0/alloy-view/pkg/state"
)

type Handler interface {
	String() string
	ServeHomePage(res http.ResponseWriter, req *http.Request)
	ServeNodesInfoPage(res http.ResponseWriter, req *http.Request)
	ServeNodeHealthIndicator(res http.ResponseWriter, req *http.Request)
	ServeSearchTargetPage(res http.ResponseWriter, req *http.Request)
	TargetSearchSubmitHandler(res http.ResponseWriter, req *http.Request)
}

type BaseHandler struct {
	logger *slog.Logger
	State  state.AppState
}
type HandlerClustered struct {
	BaseHandler
	conf *config.ClusterConfig
}

func (h *HandlerClustered) String() string {
	return fmt.Sprintf("[ handler clustered:\nconfig:\n%s ]", h.conf)
}

func NewHandler(conf config.AlloyConfig, logger *slog.Logger) (Handler, error) {
	switch mode := conf.GetMode(); mode {
	case "cluster":
		// FIX 2: Type Assertion works directly on the interface
		cfg, ok := conf.(*config.ClusterConfig)
		if !ok {
			return nil, fmt.Errorf("mode:cluster, but cfg not correctly typed")
		}

		st := state.NewClusteredAppState(cfg.SD, cfg.Mode, logger)
		st.InitState()

		return &HandlerClustered{
			BaseHandler: BaseHandler{
				logger: logger,
				State:  st,
			},
			conf: cfg,
		}, nil

	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}
}
