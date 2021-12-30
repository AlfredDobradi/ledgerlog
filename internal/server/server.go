package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/database"
	"github.com/gorilla/mux"
)

type Service struct {
	*http.Server
	Errors chan error
}

type Option func(*Service) error

func New() (*Service, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	m := mux.NewRouter()
	h := &Handler{db}

	m.HandleFunc(RouteAPISend, h.handleSend)
	m.HandleFunc(RouteAPIRegister, h.handleRegister)
	m.HandleFunc(RouteAPIPosts, h.handlePosts)

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
