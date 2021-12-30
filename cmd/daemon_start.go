package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/database"
	"github.com/AlfredDobradi/ledgerlog/internal/database/cockroach"
	"github.com/AlfredDobradi/ledgerlog/internal/server"
)

type StartCmd struct {
	IP   config.IPAddress `help:"IP address to listen on" default:"0.0.0.0"`
	Port config.Port      `help:"Port to listen on" default:"8080"`
}

func (cmd *StartCmd) Run(ctx *Context) error {
	applyDatabaseConfig()
	applyDaemonConfig(cmd.IP, cmd.Port)

	s, err := server.New()
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
			database.Close(wg)                   // nolint
			wg.Wait()
			break Loop
		case serviceErr := <-s.Errors:
			log.Printf("Received error from HTTP service: %v", serviceErr)
		}
	}

	return nil
}

func applyDatabaseConfig() {
	dbConf := config.GetSettings().Database
	database.SetDriver(dbConf.Driver)

	switch dbConf.Driver {
	// case config.DriverBadger:
	// 	badgerdb.SetDatabasePath(dbConf.Badger.Path)
	// 	badgerdb.SetValuePath(dbConf.Badger.ValuePath)
	case config.DriverCockroach:
		cockroach.SetUser(dbConf.Postgres.User)
		cockroach.SetPassword(dbConf.Postgres.Password)
		cockroach.SetHost(dbConf.Postgres.Host)
		cockroach.SetPort(dbConf.Postgres.Port)
		cockroach.SetDatabase(dbConf.Postgres.Database)
		cockroach.SetSSLMode(dbConf.Postgres.SSLMode)
		cockroach.SetSSLRootCert(dbConf.Postgres.SSLRootCert)
		cockroach.SetCluster(dbConf.Postgres.Cluster)
	}
}

func applyDaemonConfig(ipArg config.IPAddress, portArg config.Port) {
	ip := config.GetSettings().Daemon.IP
	if ip == "" || ipArg != "0.0.0.0" {
		ip = ipArg
	}
	server.SetIPAddress(ip)

	port := config.GetSettings().Daemon.Port
	if port == 0 || portArg != 8080 {
		port = portArg
	}
	server.SetPort(port)
}
