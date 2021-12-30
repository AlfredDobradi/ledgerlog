package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/database/cockroach"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"golang.org/x/crypto/ssh"
)

var (
	errInvalidDriver = fmt.Errorf("Invalid database driver specified")
)

type DB interface {
	AddPost(models.SendPostRequest) error
	GetPosts() ([]models.Post, error)
	RegisterUser(models.RegisterRequest) error
	GetUser(map[string]string) (models.User, error)
	GetPublicKey(string) (ssh.PublicKey, error)
	Close(context.Context) error
}

func GetDB() (DB, error) {
	var db DB
	var err error
	switch config.GetSettings().Database.Driver {
	// case config.DriverBadger:
	// 	db, err = badgerdb.GetConnection()
	case config.DriverCockroach:
		db, err = cockroach.GetConnection()
	default:
		err = errInvalidDriver
	}
	return db, err
}

func Close(wg *sync.WaitGroup) error {
	switch config.GetSettings().Database.Driver {
	// case config.DriverBadger:
	// 	return badgerdb.Close(wg)
	case config.DriverCockroach:
		return cockroach.Close(wg)
	default:
		return errInvalidDriver
	}
}
