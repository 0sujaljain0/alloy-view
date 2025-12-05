package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/0sujaljain0/alloy-view/pkg/config"
	"gopkg.in/yaml.v3"
)

type Server struct {
	mux    *http.ServeMux
	port   uint16
	id     string
	conf   *config.AlloyModeConfig
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

func (s *Server) configDumpHandler(res http.ResponseWriter, req *http.Request) {
	msg, err := yaml.Marshal(*s.conf)
	if err != nil {

	}
	res.Write([]byte(msg))
}

func ConfigureServer(port uint16, id string, alloyConfig *config.AlloyModeConfig, logger *slog.Logger) *Server {
	mux := http.NewServeMux()
	server := &Server{
		mux:    mux,
		port:   port,
		id:     id,
		conf:   alloyConfig,
		logger: logger,
	}

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", server.ServeHomePage)
	mux.HandleFunc("/config", server.configDumpHandler)

	return server
}
