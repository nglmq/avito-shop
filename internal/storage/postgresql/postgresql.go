package postgresql

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(ctx context.Context, dsn string) (*Repo, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error reading dsn: %w", err)
	}

	config.MaxConns = 5
	config.MinConns = 1
	config.MaxConnIdleTime = 3 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating conn pool: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting database: %w", err)
	}

	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL, 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS balances (
		    id SERIAL PRIMARY KEY,
		    username VARCHAR(255) REFERENCES users(username),
		    balance INT NOT NULL DEFAULT 1000 CHECK (balance >= 0),
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS purchases (
			id SERIAL PRIMARY KEY,
		    username VARCHAR(255) REFERENCES users(username),
    		item_name VARCHAR(255) NOT NULL,
    		amount INT NOT NULL DEFAULT 1,
    		total_price INT NOT NULL,
    		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS transactions (
		    id SERIAL PRIMARY KEY,
		    sender_username VARCHAR(255) REFERENCES users(username),
		    receiver_username VARCHAR(255) REFERENCES users(username),
		    amount INT NOT NULL,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE INDEX IF NOT EXISTS idx_transactions_sender_username ON transactions(sender_username);
		CREATE INDEX IF NOT EXISTS idx_transactions_receiver_username ON transactions(receiver_username);
	`)
	if err != nil {
		panic(err)
	}

	return &Repo{
		db: db,
	}, nil
}
