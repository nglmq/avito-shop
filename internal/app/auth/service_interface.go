package auth

import "context"

type ServiceInterface interface {
	AuthenticateUser(ctx context.Context, username, password string) (string, error)
	RegisterUser(ctx context.Context, username, password string) (string, error)
}
