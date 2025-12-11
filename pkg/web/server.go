package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/0sujaljain0/alloy-view/pkg/config"
	"github.com/0sujaljain0/alloy-view/pkg/handler"
	// "github.com/0sujaljain0/alloy-view/pkg/state"
)

type Server struct {
	mux    *http.ServeMux
	port   uint16
	id     string
	logger *slog.Logger
}

func (s *Server) String() string {
	return fmt.Sprintf("[{port: %d}-{id: %s}]", s.port, s.id)
}

func (s *Server) Start() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
	if err != nil {
		return err
	}
	return nil
}

func ConfigureServer(port uint16, id string, alloyConfig *config.AlloyConfig, logger *slog.Logger) (*Server, error) {
	mux := http.NewServeMux()

	server := &Server{
		mux:    mux,
		port:   port,
		id:     id,
		logger: logger,
	}

	hld, err := handler.NewHandler(*alloyConfig, logger)
	if err != nil  {
		return nil, err
	}

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", hld.ServeHomePage)

	return server, nil
}
