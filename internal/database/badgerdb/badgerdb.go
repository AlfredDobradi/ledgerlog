package badgerdb

import (
	"fmt"

	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	badger "github.com/dgraph-io/badger/v3"
	"golang.org/x/crypto/ssh"
)

type DB struct {
	*badger.DB
}

var connection *DB

func GetConnection(opts badger.Options) (*DB, error) {
	if connection == nil {
		db, err := badger.Open(opts)
		if err != nil {
			return nil, err
		}

		connection = &DB{
			db,
		}
	}

	return connection, nil
}

func Close() error {
	if connection != nil {
		return connection.Close()
	}
	return fmt.Errorf("Connection doesn't exist")
}

func (d *DB) GetPublicKey(email string) (ssh.PublicKey, error) {

	return nil, nil
}

func (d *DB) RegisterUser(request models.RegisterRequest) error {

	return nil
}
