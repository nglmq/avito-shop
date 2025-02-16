package history

import (
	"context"
	"github.com/nglmq/avito-shop/internal/models"
)

type InfoServiceInterface interface {
	GetInfo(ctx context.Context, username string) (models.InfoResponse, error)
}
