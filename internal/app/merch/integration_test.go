package merch_test

import (
	"context"
	"github.com/nglmq/avito-shop/internal/storage/postgresql"
	"log"
	"testing"

	"github.com/nglmq/avito-shop/internal/app/merch"
	"github.com/nglmq/avito-shop/internal/storage"
	"log/slog"
	"os"
)

func TestBuyItemIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")

	store, err := postgresql.NewRepo(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to set up test database: %v", err)
	}

	_, _ = store.SaveUser(context.Background(), "user1", "12345")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	service := merch.New(logger, store)

	tests := []struct {
		name          string
		username      string
		itemName      string
		amount        int
		expectedError error
	}{
		{
			name:          "Success",
			username:      "user1",
			itemName:      "socks",
			amount:        1,
			expectedError: nil,
		},
		{
			name:          "InsufficientFunds",
			username:      "user1",
			itemName:      "hoody",
			amount:        10,
			expectedError: storage.ErrInsufficientFunds,
		},
		{
			name:          "InvalidAmountZero",
			username:      "user1",
			itemName:      "hoody",
			amount:        0,
			expectedError: merch.ErrInvalidAmount,
		},
		{
			name:          "InvalidAmountNegative",
			username:      "user1",
			itemName:      "hoody",
			amount:        -1,
			expectedError: merch.ErrInvalidAmount,
		},
		{
			name:          "UserNotFound",
			username:      "user2",
			itemName:      "hoody",
			amount:        1,
			expectedError: storage.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.BuyItem(context.Background(), tt.username, tt.itemName, tt.amount)
			if err == nil && tt.expectedError != nil {
				t.Fatalf("expected error %v, got nil", tt.expectedError)
			}
			if err != nil && tt.expectedError == nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
