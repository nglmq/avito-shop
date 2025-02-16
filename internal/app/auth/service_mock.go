package auth

import "context"

type ServiceMock struct {
	AuthenticateUserFunc func(ctx context.Context, username, password string) (string, error)
	RegisterUserFunc     func(ctx context.Context, username, password string) (string, error)
}

func (m *ServiceMock) AuthenticateUser(ctx context.Context, username, password string) (string, error) {
	if m.AuthenticateUserFunc != nil {
		return m.AuthenticateUserFunc(ctx, username, password)
	}
	return "", nil
}

func (m *ServiceMock) RegisterUser(ctx context.Context, username, password string) (string, error) {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(ctx, username, password)
	}
	return "", nil
}
