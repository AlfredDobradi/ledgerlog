package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/database/badgerdb"
	"github.com/gorilla/mux"
)

type Service struct {
	*http.Server
	Errors chan error
}

type Option func(*Service) error

func New(bdb *badgerdb.DB, opts ...Option) (*Service, error) {
	m := mux.NewRouter()
	h := &Handler{bdb}

	m.HandleFunc(RouteAPISend, h.handleSend)
	m.HandleFunc(RouteAPIRegister, h.handleRegister)
	m.HandleFunc(RouteAPIPosts, h.handlePosts)
	m.HandleFunc(RouteDebugKeys, h.handleDebugKeys)

	s := &Service{
		Server: &http.Server{
			Addr:    "localhost:8080",
			Handler: m,
		},
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
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
