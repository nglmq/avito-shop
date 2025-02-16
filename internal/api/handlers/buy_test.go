package handlers_test

import (
	"context"
	"errors"
	"github.com/nglmq/avito-shop/internal/api/handlers"
	"github.com/nglmq/avito-shop/internal/app/merch"
	"net/http"
	"net/http/httptest"
	"testing"
)

func validBuyRequest(item string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/buy/"+item, nil)
	ctx := context.WithValue(req.Context(), "user", "validUser")
	return req.WithContext(ctx)
}

func TestHandleBuyItem(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		mockService    *merch.ServiceMock
		expectedStatus int
	}{
		{
			name:    "Success",
			request: validBuyRequest("socks"),
			mockService: &merch.ServiceMock{
				BuyItemFunc: func(ctx context.Context, username, item string, quantity int) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized",
			request:        httptest.NewRequest(http.MethodPost, "/api/buy/socks", nil),
			mockService:    &merch.ServiceMock{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:    "InternalServerError",
			request: validBuyRequest("socks"),
			mockService: &merch.ServiceMock{
				BuyItemFunc: func(ctx context.Context, username, item string, quantity int) error {
					return errors.New("internal error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:    "ItemNotFound",
			request: validBuyRequest("item"),
			mockService: &merch.ServiceMock{
				BuyItemFunc: func(ctx context.Context, username, item string, quantity int) error {
					return merch.ErrItemNotFound
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "InvalidItemParameter",
			request: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/api/buy/", nil)
				ctx := context.WithValue(req.Context(), "user", "validUser")
				return req.WithContext(ctx)
			}(),
			mockService:    &merch.ServiceMock{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.HandleBuyItem(tt.mockService)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, tt.request)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status code %v, got %v", tt.expectedStatus, status)
			}
		})
	}
}
