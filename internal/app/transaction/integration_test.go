package transaction_test

import (
	"context"
	"github.com/nglmq/avito-shop/internal/app/transaction"
	"github.com/nglmq/avito-shop/internal/storage"
	"github.com/nglmq/avito-shop/internal/storage/postgresql"
	"log"
	"testing"

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
	_, _ = store.SaveUser(context.Background(), "user2", "12345")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	service := transaction.New(logger, store)

	tests := []struct {
		name          string
		sender        string
		receiver      string
		amount        int
		expectedError error
	}{
		{
			name:          "Success#1",
			sender:        "user1",
			receiver:      "user2",
			amount:        100,
			expectedError: nil,
		},
		{
			name:          "Success#2",
			sender:        "user2",
			receiver:      "user1",
			amount:        100,
			expectedError: nil,
		},
		{
			name:          "SenderAndReceiverAreTheSame",
			sender:        "user1",
			receiver:      "user1",
			amount:        100,
			expectedError: transaction.ErrInvalidRecipient,
		},
		{
			name:          "InsufficientFunds",
			sender:        "user1",
			receiver:      "user2",
			amount:        10000,
			expectedError: transaction.ErrInsufficientBalance,
		},
		{
			name:          "InvalidAmountZero",
			sender:        "user1",
			receiver:      "user2",
			amount:        0,
			expectedError: transaction.ErrInvalidAmount,
		},
		{
			name:          "InvalidAmountNegative",
			sender:        "user1",
			receiver:      "user2",
			amount:        -1,
			expectedError: transaction.ErrInvalidAmount,
		},
		{
			name:          "SenderNotFound",
			sender:        "user3",
			receiver:      "user2",
			amount:        100,
			expectedError: storage.ErrUserNotFound,
		},
		{
			name:          "ReceiverNotFound",
			sender:        "user1",
			receiver:      "user3",
			amount:        100,
			expectedError: storage.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SendCoins(context.Background(), tt.sender, tt.receiver, tt.amount)
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
