package cockroach

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	_ssh "github.com/AlfredDobradi/ledgerlog/internal/ssh"
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

// TODO Use user model
func (c *Conn) AddPost(email string, request models.SendPostRequest) error {
	content, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return crdbpgx.ExecuteTx(context.Background(), c, pgx.TxOptions{}, func(tx pgx.Tx) error {
		user, err := c.getUserByEmail(tx, email)
		if err != nil {
			return err
		}

		prevID, err := c.getLastEntryUUID(tx)
		if err != nil {
			return err
		}

		postID := uuid.New()
		if _, err := tx.Exec(context.TODO(), `INSERT INTO snapshot_posts (id, idowner, post) VALUES($1::uuid, $2::uuid, $3)`,
			postID,
			user.ID,
			request.Message,
		); err != nil {
			return fmt.Errorf("Inserting to snapshot: %w", err)
		}

		return c.appendToLedger(tx, string(content), prevID, postID)
	})
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
		if err := c.guardEmailAndPublicKey(tx, request.Email, request.PublicKey); err != nil {
			return err
		}

		prevID, err := c.getLastEntryUUID(tx)
		if err != nil {
			return err
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

		return c.appendToLedger(tx, string(content), prevID, userID)
	})
}

func (c *Conn) GetPublicKey(email string) (ssh.PublicKey, error) {
	row := c.QueryRow(context.TODO(), "SELECT public_key FROM snapshot_users WHERE email = $1", email)
	var public_key string
	if err := row.Scan(&public_key); err != nil {
		return nil, err
	}

	return _ssh.ParsePublicKey([]byte(public_key))
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

func (c *Conn) guardEmailAndPublicKey(tx pgx.Tx, email, publicKey string) error {
	rows, err := c.Query(context.TODO(), `SELECT id FROM snapshot_users WHERE email = $1 OR public_key = $2`,
		email,
		publicKey,
	)
	if err != nil {
		return fmt.Errorf("rows: %w", err)
	}
	if rows.Next() {
		return fmt.Errorf("There is already a ledger entry for creating a user with either or both of these credentials")
	}
	return nil
}

func (c *Conn) getLastEntryUUID(tx pgx.Tx) (uuid.UUID, error) {
	last := tx.QueryRow(context.TODO(), "SELECT id FROM ledger ORDER BY t DESC LIMIT 1")
	lastID := uuid.UUID{}
	if err := last.Scan(&lastID); err != nil && err != pgx.ErrNoRows {
		return uuid.Nil, fmt.Errorf("getting last record: %w", err)
	} else if err == pgx.ErrNoRows {
		return uuid.Nil, nil
	}
	return lastID, nil
}

func (c *Conn) getUserByEmail(tx pgx.Tx, email string) (models.User, error) {
	userRow := tx.QueryRow(context.TODO(), "SELECT * FROM snapshot_users WHERE email = $1", email)
	var user models.User
	if err := userRow.Scan(&user.ID, &user.Email, &user.PreferredName, &user.PublicKey, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return models.User{}, fmt.Errorf("getting user: %w", err)
	}
	return user, nil
}

func (c *Conn) appendToLedger(tx pgx.Tx, content string, prev uuid.UUID, subject uuid.UUID) error {
	_, err := tx.Exec(context.TODO(), `INSERT INTO ledger (t, id, prev, idsubject, content) VALUES ($1, gen_random_uuid(), $2::uuid, $3::uuid, $4)`,
		time.Now(),
		prev,
		subject,
		content,
	)
	return err
}
