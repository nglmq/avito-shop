package history

import (
	"context"
	"github.com/nglmq/avito-shop/internal/models"
)

type InfoServiceMock struct {
	GetInfoFunc func(ctx context.Context, username string) (models.InfoResponse, error)
}

func (m *InfoServiceMock) GetInfo(ctx context.Context, username string) (models.InfoResponse, error) {
	if m.GetInfoFunc != nil {
		return m.GetInfoFunc(ctx, username)
	}
	return models.InfoResponse{}, nil
}
