package transaction

import "context"

type ServiceInterface interface {
	SendCoins(ctx context.Context, from, to string, amount int) error
}
