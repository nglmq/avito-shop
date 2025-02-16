package handlers

import (
	"encoding/json"
	"errors"
	"github.com/nglmq/avito-shop/internal/models"
	"log/slog"
	"net/http"
)

var (
	ErrInternal     = errors.New("internal server error")
	ErrInvalidBody  = errors.New("invalid request body")
	ErrUnauthorized = errors.New("unauthorized")
)

func respondWithError(w http.ResponseWriter, statusCode int, handlerName string, err error) {
	slog.Error("Error occurred in handler",
		slog.String("handler", handlerName),
		slog.String("error", err.Error()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := models.ErrorResponse{
		Errors: err.Error(),
	}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		slog.Error("Failed to send error response", "error", err)
	}
}
