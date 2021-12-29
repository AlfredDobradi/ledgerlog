package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/database/badgerdb"
	"github.com/AlfredDobradi/ledgerlog/internal/server"
	"github.com/dgraph-io/badger/v3"
)

type IPAddress string
type Port int

func (a IPAddress) Validate() error {
	if net.ParseIP(string(a)) == nil {
		return fmt.Errorf("Invalid IP address %s", a)
	}
	return nil
}

func (p Port) Validate() error {
	maxPort := 2 << 15
	if p <= 0 || int(p) >= maxPort {
		return fmt.Errorf("Invalid port %d", p)
	}
	return nil
}

type StartCmd struct {
	IP                IPAddress `help:"IP address to listen on" default:"0.0.0.0"`
	Port              Port      `help:"Port to listen on" default:"8080"`
	DatabasePath      string    `help:"Path to the database files" default:"./data" type:"path"`
	DatabaseValuePath string    `help:"Path to the database value files" type:"path"`
}

func (cmd *StartCmd) Run(ctx *Context) error {
	opts := badger.DefaultOptions(cmd.DatabasePath)
	if cmd.DatabaseValuePath != "" {
		opts.ValueDir = cmd.DatabaseValuePath
	}
	bdb, err := badgerdb.GetConnection(opts)
	if err != nil {
		log.Panicln(err)
	}

	s, err := server.New(bdb)
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
