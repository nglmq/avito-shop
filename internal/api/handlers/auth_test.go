package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/nglmq/avito-shop/internal/api/handlers"
	"github.com/nglmq/avito-shop/internal/app/auth"
	"github.com/nglmq/avito-shop/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func validAuthRequest() models.AuthRequest {
	return models.AuthRequest{
		Username: "validUser",
		Password: "validPass",
	}
}

func TestHandleAuth(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockService    *auth.ServiceMock
		expectedStatus int
		expectedToken  string
	}{
		{
			name:        "Success",
			requestBody: validAuthRequest(),
			mockService: &auth.ServiceMock{
				AuthenticateUserFunc: func(ctx context.Context, username, password string) (string, error) {
					return "validToken", nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedToken:  "validToken",
		},
		{
			name:           "InvalidBody",
			requestBody:    "invalid body",
			mockService:    &auth.ServiceMock{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Unauthorized",
			requestBody: validAuthRequest(),
			mockService: &auth.ServiceMock{
				AuthenticateUserFunc: func(ctx context.Context, username, password string) (string, error) {
					return "", errors.New("unauthorized")
				},
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "EmptyUsername",
			requestBody:    models.AuthRequest{Username: "", Password: "validPass"},
			mockService:    &auth.ServiceMock{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "EmptyPassword",
			requestBody:    models.AuthRequest{Username: "validUser", Password: ""},
			mockService:    &auth.ServiceMock{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.HandleAuth(tt.mockService)

			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status %v; got %v", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectedToken != "" {
				var authResp models.AuthResponse
				if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
					t.Fatalf("could not decode response: %v", err)
				}

				if authResp.Token != tt.expectedToken {
					t.Fatalf("expected token '%v'; got %v", tt.expectedToken, authResp.Token)
				}
			}
		})
	}
}
