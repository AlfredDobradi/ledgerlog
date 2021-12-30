package database

import (
	"fmt"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/database/badgerdb"
	"github.com/AlfredDobradi/ledgerlog/internal/database/cockroach"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"golang.org/x/crypto/ssh"
)

var (
	errInvalidDriver = fmt.Errorf("Invalid database driver specified")
)

type DB interface {
	AddPost(string, models.SendPostRequest) error
	GetPosts() ([]models.Post, error)
	GetKeys() ([]byte, error)
	RegisterUser(models.RegisterRequest) error
	GetPublicKey(string) (ssh.PublicKey, error)
}

func GetDB() (DB, error) {
	var db DB
	var err error
	switch config.GetSettings().Database.Driver {
	case config.DriverBadger:
		db, err = badgerdb.GetConnection(config.GetSettings().Database.Badger)
	case config.DriverCockroach:
		db, err = cockroach.GetConnection(config.GetSettings().Database.Postgres)
	default:
		err = errInvalidDriver
	}
	return db, err
}

func Close(wg *sync.WaitGroup) error {
	switch config.GetSettings().Database.Driver {
	case config.DriverBadger:
		return badgerdb.Close(wg)
	case config.DriverCockroach:
		return cockroach.Close(wg)
	default:
		return errInvalidDriver
	}
}
