package postgresql

import (
	"context"
)

func (r *Repo) CreateTransaction(ctx context.Context, senderUsername, receiverUsername string, amount int) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO transactions (sender_username, receiver_username, amount)
		VALUES ($1, $2, $3)
	`, senderUsername, receiverUsername, amount)
	if err != nil {
		return err
	}

	return nil
}
