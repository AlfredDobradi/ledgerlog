package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/database/badgerdb"
	"github.com/AlfredDobradi/ledgerlog/internal/server"
	"github.com/dgraph-io/badger/v3"
)

type StartCmd struct {
	IP                config.IPAddress `help:"IP address to listen on" default:"0.0.0.0"`
	Port              config.Port      `help:"Port to listen on" default:"8080"`
	DatabasePath      string           `help:"Path to the database files" default:"./data" type:"path"`
	DatabaseValuePath string           `help:"Path to the database value files" type:"path"`
}

func (cmd *StartCmd) Run(ctx *Context) error {
	dbPath := config.GetSettings().Database.Path
	if cmd.DatabasePath != "" {
		dbPath = cmd.DatabasePath
	}
	dbValuePath := config.GetSettings().Database.Path
	if cmd.DatabaseValuePath != "" {
		dbValuePath = cmd.DatabaseValuePath
	}

	opts := badger.DefaultOptions(dbPath)
	if dbValuePath != "" {
		opts.ValueDir = dbValuePath
	}
	bdb, err := badgerdb.GetConnection(opts)
	if err != nil {
		log.Panicln(err)
	}

	ip := config.GetSettings().Daemon.IP
	port := config.GetSettings().Daemon.Port
	if ip == "" || cmd.IP != "0.0.0.0" {
		ip = cmd.IP
	}
	if port == 0 || cmd.Port != 8080 {
		port = cmd.Port
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	s, err := server.New(bdb, server.WithAddress(addr))
	if err != nil {
		log.Panicln(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

Loop:
	for {
		select {
		case <-sigs:
			wg := &sync.WaitGroup{}
			log.Println("Received signal")
			wg.Add(2)
			s.Shutdown(context.Background(), wg) // nolint
			badgerdb.Close(wg)                   // nolint
			wg.Wait()
			break Loop
		case serviceErr := <-s.Errors:
			log.Printf("Received error from HTTP service: %v", serviceErr)
		}
	}

	return nil
}
