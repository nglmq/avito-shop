package storage

import (
	"context"
	"errors"
)

var (
	ErrUsernameExists    = errors.New("username already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type Getter interface {
	GetUser(ctx context.Context, username string) (string, string, error)
}
