package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/database/cockroach"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"golang.org/x/crypto/ssh"
)

var (
	errInvalidDriver = fmt.Errorf("Invalid database driver specified")
)

type Connection interface {
	AddPost(request models.SendPostRequest) error
	AddChannel(request models.AddChannel) error
	GetPosts(pageNum int, postsPerPage int) ([]models.PostDisplay, error)
	GetPostsSince(max int, since time.Time) ([]models.PostDisplay, int, error)
	RegisterUser(request models.RegisterRequest) error
	FindUser(filters map[string]string) (models.User, error)
	GetPublicKey(email string) (ssh.PublicKey, error)
	Close(context context.Context) error
}

func GetConnection(ctx context.Context) (Connection, error) {
	var db Connection
	var err error
	switch config.GetSettings().Database.Driver {
	// case config.DriverBadger:
	// 	db, err = badgerdb.GetConnection()
	case config.DriverCockroach:
		db, err = cockroach.GetConnection(ctx)
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
