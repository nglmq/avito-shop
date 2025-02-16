package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nglmq/avito-shop/internal/models"
	"github.com/nglmq/avito-shop/internal/storage"
)

func (r *Repo) GetInfo(ctx context.Context, username string) (models.InfoResponse, error) {
	var info models.InfoResponse

	err := r.db.QueryRow(ctx, `
		SELECT balance
		FROM balances
		WHERE username = $1
	`, username).Scan(&info.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.InfoResponse{}, storage.ErrUserNotFound
		}
		return models.InfoResponse{}, fmt.Errorf("error fetching balance: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT item_name, SUM(amount) AS total_quantity
		FROM purchases 
		WHERE username = $1
		GROUP BY item_name
	`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.InfoResponse{}, storage.ErrUserNotFound
		}
		return models.InfoResponse{}, fmt.Errorf("error fetching inventory: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item string
		var quantity int
		if err := rows.Scan(&item, &quantity); err != nil {
			return models.InfoResponse{}, fmt.Errorf("error scanning inventory row: %w", err)
		}
		info.Inventory = append(info.Inventory, models.InventoryItem{
			Item:     item,
			Quantity: quantity,
		})
	}

	transactionRows, err := r.db.Query(ctx, `
		SELECT 
			sender_username,
			receiver_username,
			amount
		FROM transactions
		WHERE sender_username = $1 OR receiver_username = $1
	`, username)
	if err != nil {
		return models.InfoResponse{}, fmt.Errorf("error fetching transaction history: %w", err)
	}
	defer transactionRows.Close()

	for transactionRows.Next() {
		var senderUsername, receiverUsername string
		var amount int
		if err := transactionRows.Scan(&senderUsername, &receiverUsername, &amount); err != nil {
			return models.InfoResponse{}, fmt.Errorf("error scanning transaction history row: %w", err)
		}

		if senderUsername == username {
			info.CoinHistory.Sent = append(info.CoinHistory.Sent, models.TransactionSentHistory{
				ToUser: receiverUsername,
				Amount: amount,
			})
		} else if receiverUsername == username {
			info.CoinHistory.Received = append(info.CoinHistory.Received, models.TransactionReceivedHistory{
				FromUser: senderUsername,
				Amount:   amount,
			})
		}
	}

	if err := transactionRows.Err(); err != nil {
		return models.InfoResponse{}, fmt.Errorf("error reading transaction rows: %w", err)
	}

	return info, nil
}
