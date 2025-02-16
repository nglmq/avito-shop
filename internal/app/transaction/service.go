package transaction

import (
	"context"
	"errors"
	"fmt"
	"github.com/nglmq/avito-shop/internal/storage"
	"log/slog"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInvalidRecipient    = errors.New("invalid recipient")
)

type Service struct {
	logger *slog.Logger
	repo   Repository
}

func New(logger *slog.Logger, repo Repository) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) SendCoins(ctx context.Context, from, to string, amount int) error {
	if from == to {
		return ErrInvalidRecipient
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}

	senderBalance, err := s.repo.GetBalance(ctx, from)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return err
		}
		return fmt.Errorf("error fetching sender balance: %w", err)
	}
	if senderBalance < amount {
		return ErrInsufficientBalance
	}

	exists, err := s.repo.GetUserByUsername(ctx, to)
	if err != nil {
		return storage.ErrUserNotFound
	}

	if !exists {
		return storage.ErrUserNotFound
	}

	err = s.repo.UpdateBalanceDeduct(ctx, from, amount)
	if err != nil {
		s.logger.Error("Error deducting balance",
			slog.String("username", from),
			slog.Int("amount", amount),
			slog.String("error", err.Error()))
		return err
	}

	err = s.repo.UpdateBalance(ctx, to, amount)
	if err != nil {
		s.logger.Error("Error updating balance",
			slog.String("username", to),
			slog.Int("amount", amount),
			slog.String("error", err.Error()))
		return err
	}

	err = s.repo.CreateTransaction(ctx, from, to, amount)
	if err != nil {
		s.logger.Error("Error creating transaction",
			slog.String("from", from),
			slog.String("to", to),
			slog.Int("amount", amount),
			slog.String("error", err.Error()))
		return err
	}

	return nil
}
