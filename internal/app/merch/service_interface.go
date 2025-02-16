package merch

import "context"

type ServiceInterface interface {
	BuyItem(ctx context.Context, username, itemName string, amount int) error
}
