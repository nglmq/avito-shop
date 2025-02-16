package auth

import "context"

type Repository interface {
	GetUserPassword(ctx context.Context, username string) (string, error)
	SaveUser(ctx context.Context, username, password string) (string, error)
}
