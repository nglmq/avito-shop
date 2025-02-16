package merch

import "context"

type Repository interface {
	GetBalance(ctx context.Context, username string) (int, error)
	UpdateBalance(ctx context.Context, receiverUUID string, amount int) error
	UpdateBalanceDeduct(ctx context.Context, senderUUID string, amount int) error
	AddPurchase(ctx context.Context, username, itemName string, amount, totalPrice int) error
}
