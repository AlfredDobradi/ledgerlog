package badgerdb

import (
	"fmt"
	"strings"

	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	_ssh "github.com/AlfredDobradi/ledgerlog/internal/ssh"
	badger "github.com/dgraph-io/badger/v3"
	"golang.org/x/crypto/ssh"
)

const (
	UserRecordKey string = "public_key:%s"
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
	var rawPubKey []byte
	err := d.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf(UserRecordKey, email)
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			rawPubKey = append([]byte{}, val...)
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return _ssh.ParsePublicKey(rawPubKey)
}

func (d *DB) RegisterUser(request models.RegisterRequest) error {
	return d.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf(UserRecordKey, request.Email)
		if _, err := txn.Get([]byte(key)); err != nil && err != badger.ErrKeyNotFound {
			return err
		} else if err == nil {
			return fmt.Errorf("user record already exists, won't update")
		}

		return txn.Set([]byte(key), []byte(request.PublicKey))
	})
}

func (d *DB) GetKeys() ([]byte, error) {
	output := make([]string, 0)
	err := d.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		iter := txn.NewIterator(opts)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			it := iter.Item()
			if err := it.Value(func(val []byte) error {
				record := fmt.Sprintf("%s > %s", it.Key(), val)
				output = append(output, record)
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return []byte(strings.Join(output, "\n")), nil
}
