package cockroach

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"github.com/cockroachdb/cockroach-go/crdb"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/ssh"
)

var errNotImplemented = fmt.Errorf("Not implemented")

type Conn struct {
	*sql.DB
}

func GetConnection() (*Conn, error) {
	connectionString := buildConnectionString()

	c, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := c.Ping(); err != nil {
		return nil, err
	}

	return &Conn{c}, nil
}

func Close(wg *sync.WaitGroup) error {
	wg.Done()
	return nil
}

func (c *Conn) AddPost(string, models.SendPostRequest) error {
	return errNotImplemented
}

func (c *Conn) GetPosts() ([]models.Post, error) {
	return nil, errNotImplemented
}

func (c *Conn) GetKeys() ([]byte, error) {
	_ = crdb.ExecuteTx(context.Background(), c.DB, nil, func(tx *sql.Tx) error {
		return nil
	})
	return nil, errNotImplemented
}

func (c *Conn) RegisterUser(models.RegisterRequest) error {
	return crdb.ExecuteTx(context.Background(), c.DB, nil, func(tx *sql.Tx) error {
		last := tx.QueryRow("SELECT id FROM ledger ORDER BY t DESC LIMIT 1")
		row := models.LedgerEntry{}
		err := last.Scan(&row)
		if err != nil {
			return last.Err()
		}

		log.Printf("%+v", row)

		return nil
	})
}

func (c *Conn) GetPublicKey(string) (ssh.PublicKey, error) {
	return nil, errNotImplemented
}

func buildConnectionString() string {
	connectionURL := url.URL{}
	connectionURL.Scheme = "postgresql"
	if User() != "" {
		if Password() != "" {
			connectionURL.User = url.UserPassword(User(), Password())
		} else {
			connectionURL.User = url.User(User())
		}
	}
	connectionURL.Host = fmt.Sprintf("%s:%s", Host(), Port())

	if Cluster() != "" {
		connectionURL.Path = fmt.Sprintf("%s.%s", Cluster(), Database())
	} else {
		connectionURL.Path = Database()
	}
	values := url.Values{}
	values.Set("sslmode", SSLMode())
	if SSLMode() != "disabled" && SSLRootCert() != "" {
		values.Set("sslrootcert", SSLRootCert())
	}
	connectionURL.RawQuery = values.Encode()

	return connectionURL.String()
}
