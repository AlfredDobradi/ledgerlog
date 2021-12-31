package cockroach

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
	"github.com/AlfredDobradi/ledgerlog/internal/server/models"
	_ssh "github.com/AlfredDobradi/ledgerlog/internal/ssh"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	"golang.org/x/crypto/ssh"
)

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
func (c *Conn) AddPost(request models.SendPostRequest) error {
	content, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return crdbpgx.ExecuteTx(context.Background(), c, pgx.TxOptions{}, func(tx pgx.Tx) error {
		prevID, err := c.getLastEntryUUID(tx)
		if err != nil {
			return err
		}

		postID := uuid.New()
		if _, err := tx.Exec(context.TODO(), `INSERT INTO snapshot_posts (id, idowner, post) VALUES($1::uuid, $2::uuid, $3)`,
			postID,
			request.Owner,
			request.Message,
		); err != nil {
			return fmt.Errorf("Inserting to snapshot: %w", err)
		}

		return c.appendToLedger(tx, string(content), prevID, postID, config.KindPost)
	})
}

func (c *Conn) GetPosts(page, num int) ([]models.PostDisplay, error) {
	if num < 1 {
		num = 0
	} else if num > 100 {
		num = 100
	}
	offset := (page - 1) * num

	posts := make([]models.PostDisplay, 0)

	err := crdbpgx.ExecuteTx(context.TODO(), c, pgx.TxOptions{}, func(tx pgx.Tx) error {
		rows, err := tx.Query(context.TODO(), "SELECT sp.id, sp.post, sp.created_at, su.id, su.preferred_name, su.public_key FROM snapshot_posts sp JOIN snapshot_users su ON sp.idowner = su.id ORDER BY created_at DESC LIMIT $1 OFFSET $2", num, offset)
		if err != nil {
			return err
		}

		var row models.PostDisplay
		var pubkey string
		for rows.Next() {
			err := rows.Scan(&row.ID, &row.Message, &row.Timestamp, &row.UserID, &row.UserPreferredName, &pubkey)
			if err != nil {
				return err
			}

			pubKey, err := _ssh.ParsePublicKey([]byte(pubkey))
			if err != nil {
				return err
			}

			row.UserFingerprint = ssh.FingerprintSHA256(pubKey)
		}
		posts = append(posts, row)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (c *Conn) GetPostsSince(max int, since time.Time) ([]models.PostDisplay, int, error) {
	posts := make([]models.PostDisplay, 0)
	rowNum := 0

	err := crdbpgx.ExecuteTx(context.TODO(), c, pgx.TxOptions{}, func(tx pgx.Tx) error {
		allCountRow := tx.QueryRow(context.TODO(), `SELECT COUNT(id) FROM snapshot_posts WHERE created_at > $1::TIMESTAMP`, since.Format(time.RFC3339Nano))
		if err := allCountRow.Scan(&rowNum); err != nil {
			return err
		}

		query := `SELECT
		sp.id, sp.post, sp.created_at, su.id, su.preferred_name, su.public_key
		FROM snapshot_posts sp
		JOIN snapshot_users su ON sp.idowner = su.id
		WHERE sp.created_at > $1::TIMESTAMP
		ORDER BY sp.created_at DESC`
		params := []interface{}{since.Format(time.RFC3339Nano)}
		if max > 0 {
			query += " LIMIT $2"
			params = append(params, max)
		}

		rows, err := tx.Query(context.TODO(), query, params...)
		if err != nil {
			return err
		}

		var row models.PostDisplay
		var pubkey string
		for rows.Next() {
			err := rows.Scan(&row.ID, &row.Message, &row.Timestamp, &row.UserID, &row.UserPreferredName, &pubkey)
			if err != nil {
				return err
			}

			pubKey, err := _ssh.ParsePublicKey([]byte(pubkey))
			if err != nil {
				return err
			}

			row.UserFingerprint = ssh.FingerprintSHA256(pubKey)
			posts = append(posts, row)
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return posts, rowNum, nil
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

		return c.appendToLedger(tx, string(content), prevID, userID, config.KindUser)
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
	rows, err := tx.Query(context.TODO(), `SELECT id FROM snapshot_users WHERE email = $1 OR public_key = $2`,
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

func (c *Conn) FindUser(filter map[string]string) (models.User, error) {
	validFilter := map[string]string{}
	for key, value := range filter {
		switch key {
		case "id", "email", "preferred_name", "public_key":
			validFilter[key] = value
		}
	}
	filterStr := []string{}
	values := []interface{}{}
	index := 1
	for key, value := range validFilter {
		filterStr = append(filterStr, fmt.Sprintf("%s = $%d", key, index))
		values = append(values, value)
		index++
	}
	query := fmt.Sprintf("SELECT * FROM snapshot_users WHERE %s", strings.Join(filterStr, " AND "))
	userRow := c.QueryRow(context.TODO(), query, values...)
	var user models.User
	var publicKey string
	if err := userRow.Scan(&user.ID, &user.Email, &user.PreferredName, &publicKey, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return models.User{}, fmt.Errorf("getting user: %w", err)
	}

	var pkeyErr error
	user.PublicKey, pkeyErr = _ssh.ParsePublicKey([]byte(publicKey))
	if pkeyErr != nil {
		return models.User{}, pkeyErr
	}

	return user, nil
}

func (c *Conn) appendToLedger(tx pgx.Tx, content string, prev uuid.UUID, subject uuid.UUID, kind config.SubjectKind) error {
	_, err := tx.Exec(context.TODO(), `INSERT INTO ledger (t, id, prev, idsubject, subject_type, content) VALUES ($1, gen_random_uuid(), $2::uuid, $3::uuid, $4, $5)`,
		time.Now(),
		prev,
		subject,
		string(kind),
		content,
	)
	return err
}
