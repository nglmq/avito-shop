package transaction_test

import (
	"context"
	"errors"
	"github.com/nglmq/avito-shop/internal/app/transaction"
	"github.com/nglmq/avito-shop/internal/storage"
	"testing"
)

type MockTransactionRepository struct {
	GetBalanceFunc          func(ctx context.Context, username string) (int, error)
	GetUserByUsernameFunc   func(ctx context.Context, username string) (bool, error)
	UpdateBalanceDeductFunc func(ctx context.Context, username string, amount int) error
	UpdateBalanceFunc       func(ctx context.Context, username string, amount int) error
	CreateTransactionFunc   func(ctx context.Context, from, to string, amount int) error
}

func (m *MockTransactionRepository) GetBalance(ctx context.Context, username string) (int, error) {
	if m.GetBalanceFunc != nil {
		return m.GetBalanceFunc(ctx, username)
	}
	return 0, nil
}

func (m *MockTransactionRepository) GetUserByUsername(ctx context.Context, username string) (bool, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(ctx, username)
	}
	return false, nil
}

func (m *MockTransactionRepository) UpdateBalanceDeduct(ctx context.Context, username string, amount int) error {
	if m.UpdateBalanceDeductFunc != nil {
		return m.UpdateBalanceDeductFunc(ctx, username, amount)
	}
	return nil
}

func (m *MockTransactionRepository) UpdateBalance(ctx context.Context, username string, amount int) error {
	if m.UpdateBalanceFunc != nil {
		return m.UpdateBalanceFunc(ctx, username, amount)
	}
	return nil
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, from, to string, amount int) error {
	if m.CreateTransactionFunc != nil {
		return m.CreateTransactionFunc(ctx, from, to, amount)
	}
	return nil
}

func TestSendCoins(t *testing.T) {
	tests := []struct {
		name          string
		from          string
		to            string
		amount        int
		mockRepo      *MockTransactionRepository
		expectedError error
	}{
		{
			name:   "Success",
			from:   "user1",
			to:     "user2",
			amount: 100,
			mockRepo: &MockTransactionRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				GetUserByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
					return true, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				UpdateBalanceFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				CreateTransactionFunc: func(ctx context.Context, from, to string, amount int) error {
					return nil
				},
			},
			expectedError: nil,
		},
		{
			name:   "InsufficientBalance",
			from:   "user1",
			to:     "user2",
			amount: 100,
			mockRepo: &MockTransactionRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 50, nil
				},
				GetUserByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
					return true, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				UpdateBalanceFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				CreateTransactionFunc: func(ctx context.Context, from, to string, amount int) error {
					return nil
				},
			},
			expectedError: transaction.ErrInsufficientBalance,
		},
		{
			name:   "UserNotFound",
			from:   "user1",
			to:     "user2",
			amount: 100,
			mockRepo: &MockTransactionRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				GetUserByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
					return false, storage.ErrUserNotFound
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				UpdateBalanceFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				CreateTransactionFunc: func(ctx context.Context, from, to string, amount int) error {
					return nil
				},
			},
			expectedError: storage.ErrUserNotFound,
		},
		{
			name:   "ErrorFetchingBalance",
			from:   "user1",
			to:     "user2",
			amount: 100,
			mockRepo: &MockTransactionRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 0, errors.New("database error")
				},
				GetUserByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
					return true, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				UpdateBalanceFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				CreateTransactionFunc: func(ctx context.Context, from, to string, amount int) error {
					return nil
				},
			},
			expectedError: errors.New("error fetching sender balance: database error"),
		},
		{
			name:   "InvalidAmountZero",
			from:   "user1",
			to:     "user2",
			amount: 0,
			mockRepo: &MockTransactionRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				GetUserByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
					return true, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				UpdateBalanceFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				CreateTransactionFunc: func(ctx context.Context, from, to string, amount int) error {
					return nil
				},
			},
			expectedError: transaction.ErrInvalidAmount,
		},
		{
			name:   "InvalidAmountNegative",
			from:   "user1",
			to:     "user2",
			amount: -10,
			mockRepo: &MockTransactionRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				GetUserByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
					return true, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				UpdateBalanceFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				CreateTransactionFunc: func(ctx context.Context, from, to string, amount int) error {
					return nil
				},
			},
			expectedError: transaction.ErrInvalidAmount,
		},
		{
			name:   "SenderAndReceiverSame",
			from:   "user1",
			to:     "user1",
			amount: 100,
			mockRepo: &MockTransactionRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				GetUserByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
					return true, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				UpdateBalanceFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				CreateTransactionFunc: func(ctx context.Context, from, to string, amount int) error {
					return nil
				},
			},
			expectedError: transaction.ErrInvalidRecipient,
		},
		{
			name:   "NegativeBalanceAfterTransaction",
			from:   "user1",
			to:     "user2",
			amount: 1100,
			mockRepo: &MockTransactionRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				GetUserByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
					return true, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				UpdateBalanceFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				CreateTransactionFunc: func(ctx context.Context, from, to string, amount int) error {
					return nil
				},
			},
			expectedError: transaction.ErrInsufficientBalance,
		},
		{
			name:   "ZeroBalanceAfterTransaction",
			from:   "user1",
			to:     "user2",
			amount: 1000,
			mockRepo: &MockTransactionRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				GetUserByUsernameFunc: func(ctx context.Context, username string) (bool, error) {
					return true, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				UpdateBalanceFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				CreateTransactionFunc: func(ctx context.Context, from, to string, amount int) error {
					return nil
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := transaction.New(nil, tt.mockRepo)
			err := service.SendCoins(context.Background(), tt.from, tt.to, tt.amount)
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
