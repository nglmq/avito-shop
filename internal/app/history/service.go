package history

import (
	"context"
	"github.com/nglmq/avito-shop/internal/models"
	"log/slog"
)

type InfoService struct {
	repository InfoRepository
	logger     *slog.Logger
}

func New(logger *slog.Logger, repository InfoRepository) *InfoService {
	return &InfoService{
		logger:     logger,
		repository: repository,
	}
}

func (s *InfoService) GetInfo(ctx context.Context, username string) (models.InfoResponse, error) {
	info, err := s.repository.GetInfo(ctx, username)
	if err != nil {
		s.logger.Error("Error getting user info",
			slog.String("username", username),
			slog.String("error", err.Error()))
		return models.InfoResponse{}, err
	}

	return info, nil
}
