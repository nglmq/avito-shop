package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/nglmq/avito-shop/internal/storage"
)

func (r *Repo) GetBalance(ctx context.Context, username string) (int, error) {
	var balance int

	err := r.db.QueryRow(ctx, `
		SELECT balance FROM balances WHERE username = $1 FOR UPDATE
	`, username).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrUserNotFound
		}
		return 0, err
	}

	return balance, nil
}

func (r *Repo) UpdateBalanceDeduct(ctx context.Context, senderUsername string, amount int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE balances SET balance = balance - $1 WHERE username = $2", amount, senderUsername)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) UpdateBalance(ctx context.Context, receiverUsername string, amount int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE balances SET balance = balance + $1 WHERE username = $2", amount, receiverUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrUserNotFound
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
