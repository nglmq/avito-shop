package merch_test

import (
	"context"
	"errors"
	"github.com/nglmq/avito-shop/internal/app/merch"
	"github.com/nglmq/avito-shop/internal/storage"
	"testing"
)

type MockMerchRepository struct {
	GetBalanceFunc          func(ctx context.Context, username string) (int, error)
	UpdateBalanceDeductFunc func(ctx context.Context, username string, amount int) error
	AddPurchaseFunc         func(ctx context.Context, username, itemName string, amount, totalPrice int) error
	UpdateBalanceFunc       func(ctx context.Context, receiverUUID string, amount int) error
}

func (m *MockMerchRepository) GetBalance(ctx context.Context, username string) (int, error) {
	return m.GetBalanceFunc(ctx, username)
}

func (m *MockMerchRepository) UpdateBalanceDeduct(ctx context.Context, username string, amount int) error {
	return m.UpdateBalanceDeductFunc(ctx, username, amount)
}

func (m *MockMerchRepository) AddPurchase(ctx context.Context, username, itemName string, amount, totalPrice int) error {
	return m.AddPurchaseFunc(ctx, username, itemName, amount, totalPrice)
}

func (m *MockMerchRepository) UpdateBalance(ctx context.Context, receiverUUID string, amount int) error {
	return m.UpdateBalanceFunc(ctx, receiverUUID, amount)
}

func TestBuyItem(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		itemName      string
		amount        int
		mockRepo      *MockMerchRepository
		expectedError error
	}{
		{
			name:     "Success",
			username: "user1",
			itemName: "socks",
			amount:   1,
			mockRepo: &MockMerchRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				AddPurchaseFunc: func(ctx context.Context, username, itemName string, amount, totalPrice int) error {
					return nil
				},
			},
			expectedError: nil,
		},
		{
			name:     "InvalidAmountZero",
			username: "user1",
			itemName: "socks",
			amount:   0,
			mockRepo: &MockMerchRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				AddPurchaseFunc: func(ctx context.Context, username, itemName string, amount, totalPrice int) error {
					return nil
				},
			},
			expectedError: merch.ErrInvalidAmount,
		},
		{
			name:     "InvalidAmountNegative",
			username: "user1",
			itemName: "socks",
			amount:   -1,
			mockRepo: &MockMerchRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				AddPurchaseFunc: func(ctx context.Context, username, itemName string, amount, totalPrice int) error {
					return nil
				},
			},
			expectedError: merch.ErrInvalidAmount,
		},
		{
			name:     "InsufficientFunds",
			username: "user1",
			itemName: "hoody",
			amount:   1,
			mockRepo: &MockMerchRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 100, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				AddPurchaseFunc: func(ctx context.Context, username, itemName string, amount, totalPrice int) error {
					return nil
				},
			},
			expectedError: merch.ErrInsufficientBalance,
		},
		{
			name:     "ItemNotFound",
			username: "user1",
			itemName: "nonexistentItem",
			amount:   1,
			mockRepo: &MockMerchRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				AddPurchaseFunc: func(ctx context.Context, username, itemName string, amount, totalPrice int) error {
					return nil
				},
			},
			expectedError: merch.ErrItemNotFound,
		},
		{
			name:     "UserNotFound",
			username: "nonexistentUser",
			itemName: "socks",
			amount:   1,
			mockRepo: &MockMerchRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 0, storage.ErrUserNotFound
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				AddPurchaseFunc: func(ctx context.Context, username, itemName string, amount, totalPrice int) error {
					return nil
				},
			},
			expectedError: storage.ErrUserNotFound,
		},
		{
			name:     "ErrorFetchingBalance",
			username: "user1",
			itemName: "socks",
			amount:   1,
			mockRepo: &MockMerchRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 0, errors.New("database error")
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				AddPurchaseFunc: func(ctx context.Context, username, itemName string, amount, totalPrice int) error {
					return nil
				},
			},
			expectedError: errors.New("error fetching user balance: database error"),
		},
		{
			name:     "ErrorDeductingBalance",
			username: "user1",
			itemName: "socks",
			amount:   1,
			mockRepo: &MockMerchRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return errors.New("deduction error")
				},
				AddPurchaseFunc: func(ctx context.Context, username, itemName string, amount, totalPrice int) error {
					return nil
				},
			},
			expectedError: errors.New("error deducting balance: deduction error"),
		},
		{
			name:     "ErrorAddingPurchase",
			username: "user1",
			itemName: "socks",
			amount:   1,
			mockRepo: &MockMerchRepository{
				GetBalanceFunc: func(ctx context.Context, username string) (int, error) {
					return 1000, nil
				},
				UpdateBalanceDeductFunc: func(ctx context.Context, username string, amount int) error {
					return nil
				},
				AddPurchaseFunc: func(ctx context.Context, username, itemName string, amount, totalPrice int) error {
					return errors.New("purchase error")
				},
			},
			expectedError: errors.New("error adding purchase: purchase error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := merch.New(nil, tt.mockRepo)
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
