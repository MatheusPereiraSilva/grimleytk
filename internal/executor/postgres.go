package executor

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/MatheusPereiraSilva/grimleytk/internal/config"
	"github.com/MatheusPereiraSilva/grimleytk/internal/planner"
)

// PostgresExecutor executes plans against a PostgreSQL database
type PostgresExecutor struct {
	db *sql.DB
}

// NewPostgresExecutor creates a new Postgres executor
func NewPostgresExecutor(cfg config.Database) (*PostgresExecutor, error) {
	password := os.Getenv(cfg.Credentials.PasswordEnv)
	if password == "" {
		return nil, fmt.Errorf(
			"environment variable %s is not set",
			cfg.Credentials.PasswordEnv,
		)
	}

	conn := &url.URL{Scheme: "postgres", User: url.UserPassword(cfg.Credentials.User, password), Host: net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)), Path: "/" + cfg.Name}
	query := conn.Query()
	query.Set("sslmode", sslMode(cfg.SSL))
	query.Set("connect_timeout", "10")
	conn.RawQuery = query.Encode()
	connStr := conn.String()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Validate connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &PostgresExecutor{db: db}, nil
}

// Execute executes all actions inside a single transaction
func (e *PostgresExecutor) Execute(
	ctx context.Context,
	actions []planner.Action,
) error {

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, action := range actions {
		if _, err := tx.ExecContext(ctx, action.SQL); err != nil {
			_ = tx.Rollback()
			return &ExecutionError{
				Action: action,
				Err:    err,
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func sslMode(enabled bool) string {
	if enabled {
		return "require"
	}
	return "disable"
}

// Close releases database resources.
func (e *PostgresExecutor) Close() error { return e.db.Close() }
