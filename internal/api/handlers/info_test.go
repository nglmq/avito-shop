package handlers_test

import (
	"context"
	"errors"
	"github.com/nglmq/avito-shop/internal/api/handlers"
	"github.com/nglmq/avito-shop/internal/app/history"
	"github.com/nglmq/avito-shop/internal/models"
	"github.com/nglmq/avito-shop/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func validInfoRequest() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	ctx := context.WithValue(req.Context(), "user", "validUser")
	return req.WithContext(ctx)
}

func TestHandleGetInfo(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		mockService    *history.InfoServiceMock
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Success",
			request: validInfoRequest(),
			mockService: &history.InfoServiceMock{
				GetInfoFunc: func(ctx context.Context, username string) (models.InfoResponse, error) {
					return models.InfoResponse{}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"coins":0,"inventory":null,"coinHistory":{"received":null,"sent":null}}`,
		},
		{
			name:           "Unauthorized",
			request:        httptest.NewRequest(http.MethodGet, "/info", nil),
			mockService:    &history.InfoServiceMock{},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"unauthorized"}`,
		},
		{
			name:    "InternalServerError",
			request: validInfoRequest(),
			mockService: &history.InfoServiceMock{
				GetInfoFunc: func(ctx context.Context, username string) (models.InfoResponse, error) {
					return models.InfoResponse{}, errors.New("internal error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"internal server error"}`,
		},
		{
			name:    "UserNotFound",
			request: validInfoRequest(),
			mockService: &history.InfoServiceMock{
				GetInfoFunc: func(ctx context.Context, username string) (models.InfoResponse, error) {
					return models.InfoResponse{}, storage.ErrUserNotFound
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"internal server error"}`,
		},
		{
			name: "InvalidUserContext",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/info", nil)
				return req
			}(),
			mockService:    &history.InfoServiceMock{},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"unauthorized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.HandleGetInfo(tt.mockService)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, tt.request)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status code %v, got %v", tt.expectedStatus, status)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("expected body %v, got %v", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
