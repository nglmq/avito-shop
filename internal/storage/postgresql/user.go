package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/nglmq/avito-shop/internal/storage"
)

func (r *Repo) SaveUser(ctx context.Context, username, password string) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, password)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == pgerrcode.UniqueViolation {
			return "", storage.ErrUsernameExists
		}
		return "", fmt.Errorf("failed to insert history: %w", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO balances (username, balance) VALUES ($1, $2)", username, 1000)
	if err != nil {
		return "", fmt.Errorf("failed to insert coin balance: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return username, nil
}

func (r *Repo) GetUserPassword(ctx context.Context, username string) (string, error) {
	var userPassword string

	err := r.db.QueryRow(
		ctx,
		"SELECT password_hash FROM users WHERE username = $1", username).
		Scan(&userPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUserNotFound
		}

		return "", err
	}

	return userPassword, nil
}

func (r *Repo) GetUserByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		"SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)", username).
		Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
