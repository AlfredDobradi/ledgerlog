package cockroach

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"sync"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"github.com/cockroachdb/cockroach-go/crdb"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/ssh"
)

var errNotImplemented = fmt.Errorf("Not implemented")

type Conn struct {
	*sql.DB
}

func GetConnection(cfg config.PostgresSettings) (*Conn, error) {
	c, err := sql.Open("postgres", buildConnectionString(cfg))
	if err != nil {
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
	return errNotImplemented
}

func (c *Conn) GetPublicKey(string) (ssh.PublicKey, error) {
	return nil, errNotImplemented
}

func buildConnectionString(cfg config.PostgresSettings) string {
	connectionURL := url.URL{}
	connectionURL.Scheme = "postgresql"
	if cfg.User != "" {
		if cfg.Password != "" {
			connectionURL.User = url.UserPassword(cfg.User, cfg.Password)
		} else {
			connectionURL.User = url.User(cfg.User)
		}
	}
	connectionURL.Host = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	connectionURL.Path = cfg.Database
	values := url.Values{}
	values.Set("sslmode", cfg.SSLMode)
	if cfg.SSLMode != "disabled" && cfg.SSLRootCert != "" {
		values.Set("sslrootcert", cfg.SSLRootCert)
	}
	if cfg.Options != "" {
		values.Set("options", cfg.Options)
	}
	connectionURL.RawQuery = values.Encode()

	return connectionURL.String()
}
