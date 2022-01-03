package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service struct {
	*http.Server
	Errors chan error
}

type Option func(*Service) error

func New() (*Service, error) {
	m := mux.NewRouter()
	h, err := NewHandler()
	if err != nil {
		return nil, err
	}

	m.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			Registry:          prometheus.DefaultRegisterer,
			EnableOpenMetrics: false,
		},
	))
	m.HandleFunc(RouteAPISend, h.handleAPISend)
	m.HandleFunc(RouteAPIRegister, h.handleAPIRegister)
	m.HandleFunc(RouteAPIPosts, h.handleAPIPosts)
	m.HandleFunc(RouteAPIWS, h.handlePostsSocket)
	m.HandleFunc(RouteIndex, h.handleIndex)

	staticPath := fmt.Sprintf("%s/static", PublicPath())
	m.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	s := &Service{
		Server: &http.Server{
			Addr:    "localhost:8080",
			Handler: m,
		},
	}

	go func() {
		if config.GetSettings().Debug {
			log.Printf("Starting daemon, listening on %s", s.Addr)
		}
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Errors <- err
		}
	}()

	return s, nil
}

func (s *Service) Shutdown(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := s.Server.Shutdown(ctx)
	return err
}
