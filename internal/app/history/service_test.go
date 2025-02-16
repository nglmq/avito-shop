package history_test

import (
	"context"
	"errors"
	"github.com/nglmq/avito-shop/internal/app/history"
	"github.com/nglmq/avito-shop/internal/models"
	"github.com/nglmq/avito-shop/internal/storage"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

type MockInfoRepository struct {
	GetInfoFunc func(ctx context.Context, userUUID string) (models.InfoResponse, error)
}

func (m *MockInfoRepository) GetInfo(ctx context.Context, userUUID string) (models.InfoResponse, error) {
	if m.GetInfoFunc != nil {
		return m.GetInfoFunc(ctx, userUUID)
	}
	return models.InfoResponse{}, nil
}

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockRepo      *MockInfoRepository
		expectedInfo  models.InfoResponse
		expectedError error
	}{
		{
			name:     "Success",
			username: "user-123",
			mockRepo: &MockInfoRepository{
				GetInfoFunc: func(ctx context.Context, username string) (models.InfoResponse, error) {
					return models.InfoResponse{
						Coins: 100,
						Inventory: []models.InventoryItem{
							{Item: "socks", Quantity: 2},
						},
						CoinHistory: models.CoinHistory{
							Received: []models.TransactionReceivedHistory{
								{FromUser: "user-456", Amount: 500},
							},
							Sent: []models.TransactionSentHistory{
								{ToUser: "user-789", Amount: 300},
							},
						},
					}, nil
				},
			},
			expectedInfo: models.InfoResponse{
				Coins: 100,
				Inventory: []models.InventoryItem{
					{Item: "socks", Quantity: 2},
				},
				CoinHistory: models.CoinHistory{
					Received: []models.TransactionReceivedHistory{
						{FromUser: "user-456", Amount: 500},
					},
					Sent: []models.TransactionSentHistory{
						{ToUser: "user-789", Amount: 300},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "RepositoryError",
			username: "user-123",
			mockRepo: &MockInfoRepository{
				GetInfoFunc: func(ctx context.Context, username string) (models.InfoResponse, error) {
					return models.InfoResponse{}, errors.New("repository error")
				},
			},
			expectedInfo:  models.InfoResponse{},
			expectedError: errors.New("repository error"),
		},
		{
			name:     "InvalidUsername",
			username: "invalidUser",
			mockRepo: &MockInfoRepository{
				GetInfoFunc: func(ctx context.Context, username string) (models.InfoResponse, error) {
					return models.InfoResponse{}, storage.ErrUserNotFound
				},
			},
			expectedInfo:  models.InfoResponse{},
			expectedError: storage.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			service := history.New(logger, tt.mockRepo)

			result, err := service.GetInfo(context.Background(), tt.username)
			if err == nil && tt.expectedError != nil {
				t.Fatalf("expected error %v, got nil", tt.expectedError)
			}
			if err != nil && tt.expectedError == nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}
			if !assert.Equal(t, tt.expectedInfo, result) {
				t.Fatalf("expected info %v, got %v", tt.expectedInfo, result)
			}
		})
	}
}
