package cockroach

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	"golang.org/x/crypto/ssh"
)

var errNotImplemented = fmt.Errorf("Not implemented")

type Conn struct {
	*pgx.Conn
}

func GetConnection() (*Conn, error) {
	connectionString := buildConnectionString()

	c, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	if err := c.Ping(context.TODO()); err != nil {
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
	_ = crdbpgx.ExecuteTx(context.Background(), c, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return nil
	})
	return nil, errNotImplemented
}

func (c *Conn) RegisterUser(request models.RegisterRequest) error {
	content, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return crdbpgx.ExecuteTx(context.Background(), c, pgx.TxOptions{}, func(tx pgx.Tx) error {
		rows, err := tx.Query(context.TODO(), `SELECT id FROM ledger WHERE content @> $1::json OR content @> $2::json`,
			fmt.Sprintf(`{"email": "%s"}`, request.Email),
			fmt.Sprintf(`{"public_key": "%s"}`, request.PublicKey),
		)
		if err != nil {
			return fmt.Errorf("rows: %w", err)
		}
		if rows.Next() {
			return fmt.Errorf("There is already a ledger entry for creating a user with either or both of these credentials")
		}

		last := tx.QueryRow(context.TODO(), "SELECT id FROM ledger ORDER BY t DESC LIMIT 1")
		lastRow := models.LedgerEntry{}
		if err := last.Scan(&lastRow); err != nil && err != pgx.ErrNoRows {
			log.Printf("%T", err)
			return fmt.Errorf("getting last record: %w", err)
		}

		userID := uuid.New()
		if _, err := tx.Exec(context.TODO(), `INSERT INTO snapshot_users (id, email, preferred_name, public_key) VALUES($1::uuid, $2, $3, $4)`,
			userID,
			request.Email,
			request.Name,
			request.PublicKey,
		); err != nil {
			return fmt.Errorf("Inserting to snapshot: %w", err)
		}

		currentRow := models.LedgerEntry{
			Timestamp: time.Now(),
			ID:        uuid.New(),
			Content:   content,
			Prev:      uuid.Nil,
			Subject:   userID,
		}
		if err == nil {
			currentRow.Prev = lastRow.ID
		}

		if _, err := tx.Exec(context.TODO(), `INSERT INTO ledger (t, id, prev, idsubject, content) VALUES ($1, $2, $3::uuid, $4::uuid, $5)`,
			currentRow.Timestamp,
			currentRow.ID,
			currentRow.Prev,
			currentRow.Subject,
			string(currentRow.Content),
		); err != nil {
			return fmt.Errorf("Inserting to ledger: %w", err)
		}

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
