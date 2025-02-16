package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/nglmq/avito-shop/internal/storage"
	"github.com/nglmq/avito-shop/internal/utils/jwt"
	"github.com/nglmq/avito-shop/internal/utils/validation"
	"log/slog"
)

type Service struct {
	userRepo Repository
	logger   *slog.Logger
}

func New(logger *slog.Logger, userRepo Repository) *Service {
	return &Service{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (s *Service) AuthenticateUser(ctx context.Context, username, password string) (string, error) {
	storedPassHash, err := s.userRepo.GetUserPassword(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			s.logger.Info("User not found, proceeding with registration",
				slog.String("username", username))

			username, err = s.RegisterUser(ctx, username, password)

			if err != nil {
				s.logger.Error("Error registering user",
					slog.String("username", username),
					slog.String("error", err.Error()))
				return "", fmt.Errorf("error registering user: %w", err)
			}

			return username, nil
		}

		s.logger.Error("Error fetching user data",
			slog.String("username", username),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("error getting user password: %w", err)
	}

	if !validation.CheckPassword(password, storedPassHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := ujwt.BuildJWTString(username)
	if err != nil {
		s.logger.Error("Error generating JWT token",
			slog.String("userUUID", username),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("error generating JWT token: %w", err)
	}

	return token, nil
}

func (s *Service) RegisterUser(ctx context.Context, username, password string) (string, error) {
	passHash, err := validation.HashPassword(password)
	if err != nil {
		s.logger.Error("Error hashing password",
			slog.String("username", username),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	username, err = s.userRepo.SaveUser(ctx, username, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUsernameExists) {
			return "", err
		}

		s.logger.Error("Error saving user",
			slog.String("username", username),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("error saving user: %w", err)
	}

	token, err := ujwt.BuildJWTString(username)
	if err != nil {
		s.logger.Error("Error generating JWT token",
			slog.String("userUUID", username),
			slog.String("error", err.Error()))
		return "", fmt.Errorf("error generating JWT token: %w", err)
	}

	return token, nil
}
