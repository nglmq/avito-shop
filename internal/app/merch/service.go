package merch

import (
	"context"
	"errors"
	"fmt"
	"github.com/nglmq/avito-shop/internal/models"
	"github.com/nglmq/avito-shop/internal/storage"
	"log/slog"
)

var (
	ErrItemNotFound        = errors.New("item not found")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInsufficientBalance = errors.New("insufficient balance")
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

func (s *Service) BuyItem(ctx context.Context, username, itemName string, amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	price, exists := models.GetItemPrice(itemName)
	if !exists {
		return ErrItemNotFound
	}

	totalPrice := price * amount

	balance, err := s.repo.GetBalance(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return err
		}
		return fmt.Errorf("error fetching user balance: %w", err)
	}

	if balance < totalPrice {
		return ErrInsufficientBalance
	}

	if err := s.repo.UpdateBalanceDeduct(ctx, username, totalPrice); err != nil {
		return fmt.Errorf("error deducting balance: %w", err)
	}

	if err := s.repo.AddPurchase(ctx, username, itemName, amount, totalPrice); err != nil {
		return fmt.Errorf("error adding purchase: %w", err)
	}

	return nil
}
