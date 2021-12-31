package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/gorilla/mux"
)

type Service struct {
	*http.Server
	Errors chan error
}

type Option func(*Service) error

func New() *Service {
	m := mux.NewRouter()
	h := &Handler{}

	m.HandleFunc(RouteAPISend, h.handleAPISend)
	m.HandleFunc(RouteAPIRegister, h.handleAPIRegister)
	m.HandleFunc(RouteAPIPosts, h.handleAPIPosts)
	m.HandleFunc(RouteIndex, h.handleIndex)

	m.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static"))))

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

	return s
}

func (s *Service) Shutdown(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := s.Server.Shutdown(ctx)
	return err
}
