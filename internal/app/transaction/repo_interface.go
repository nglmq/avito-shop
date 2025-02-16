package transaction

import "context"

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (bool, error)
	CreateTransaction(ctx context.Context, senderUUID, receiverUUID string, amount int) error
	GetBalance(ctx context.Context, username string) (int, error)
	UpdateBalance(ctx context.Context, receiverUUID string, amount int) error
	UpdateBalanceDeduct(ctx context.Context, senderUUID string, amount int) error
}
