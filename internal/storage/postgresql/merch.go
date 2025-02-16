package postgresql

import (
	"context"
	"fmt"
)

func (r *Repo) AddPurchase(ctx context.Context, username, itemName string, amount, totalPrice int) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO purchases (username, item_name, amount, total_price)
		VALUES ($1, $2, $3, $4)
	`, username, itemName, amount, totalPrice)
	if err != nil {
		return fmt.Errorf("error adding purchase: %w", err)
	}
	return nil
}
