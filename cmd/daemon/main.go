package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/database/badgerdb"
	"github.com/AlfredDobradi/ledgerlog/internal/server"
	"github.com/dgraph-io/badger/v3"
)

func main() {
	opts := badger.DefaultOptions("./tmp")
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

	for {
		<-sigs
		wg := &sync.WaitGroup{}
		log.Println("Received signal")
		wg.Add(2)
		s.Shutdown(context.Background(), wg) // nolint
		badgerdb.Close(wg)                   // nolint

		wg.Wait()
		break
	}
}
