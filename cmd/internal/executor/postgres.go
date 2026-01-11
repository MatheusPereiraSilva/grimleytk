package executor

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"grimleytk/internal/config"
	"grimleytk/internal/planner"
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

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Credentials.User,
		password,
		cfg.Name,
		sslMode(cfg.SSL),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Validate connection
	if err := db.Ping(); err != nil {
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
