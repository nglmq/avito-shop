package merch

import "context"

type ServiceMock struct {
	BuyItemFunc func(ctx context.Context, username, item string, quantity int) error
}

func (m *ServiceMock) BuyItem(ctx context.Context, username, item string, quantity int) error {
	if m.BuyItemFunc != nil {
		return m.BuyItemFunc(ctx, username, item, quantity)
	}
	return nil
}
