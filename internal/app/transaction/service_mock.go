package transaction

import "context"

type ServiceMock struct {
	SendCoinsFunc func(ctx context.Context, fromUser, toUser string, amount int) error
}

func (m *ServiceMock) SendCoins(ctx context.Context, fromUser, toUser string, amount int) error {
	if m.SendCoinsFunc != nil {
		return m.SendCoinsFunc(ctx, fromUser, toUser, amount)
	}
	return nil
}
