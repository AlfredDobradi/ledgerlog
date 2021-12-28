package server

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/database/badgerdb"
	"github.com/gorilla/mux"
)

type Service struct {
	*http.Server
	stop   chan struct{}
	errors chan error
}

type Opt func(*Service) error

func New(bdb *badgerdb.DB, opts ...Opt) (*Service, error) {
	m := mux.NewRouter()
	h := &Handler{bdb}

	m.HandleFunc("/test", h.handleTest)
	m.HandleFunc("/register", h.handleRegister)
	m.HandleFunc("/debug/keys", h.handleDebugKeys)

	s := &Service{
		Server: &http.Server{
			Addr:    "localhost:8080",
			Handler: m,
		},
		stop: make(chan struct{}),
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errors <- err
		}
	}()

	return s, nil
}

func (s *Service) Shutdown(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := s.Server.Shutdown(ctx)
	return err
}
