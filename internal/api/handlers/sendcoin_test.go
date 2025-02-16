package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/nglmq/avito-shop/internal/api/handlers"
	"github.com/nglmq/avito-shop/internal/app/transaction"
	"github.com/nglmq/avito-shop/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func validSendCoinRequest() *http.Request {
	reqBody, _ := json.Marshal(models.SendCoinsRequest{ToUser: "recipient", Amount: 100})
	req := httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer(reqBody))
	ctx := context.WithValue(req.Context(), "user", "validUser")
	return req.WithContext(ctx)
}

func invalidBodyReq() *http.Request {
	reqBody, _ := json.Marshal(models.SendCoinsRequest{ToUser: "recipient"})
	req := httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer(reqBody))
	ctx := context.WithValue(req.Context(), "user", "validUser")
	return req.WithContext(ctx)
}

func TestHandleSendCoin(t *testing.T) {
	tests := []struct {
		name           string
		request        *http.Request
		mockService    *transaction.ServiceMock
		expectedStatus int
	}{
		{
			name:    "Success",
			request: validSendCoinRequest(),
			mockService: &transaction.ServiceMock{
				SendCoinsFunc: func(ctx context.Context, fromUser, toUser string, amount int) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "InvalidBody",
			request:        invalidBodyReq(),
			mockService:    &transaction.ServiceMock{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "InternalServerError",
			request: validSendCoinRequest(),
			mockService: &transaction.ServiceMock{
				SendCoinsFunc: func(ctx context.Context, fromUser, toUser string, amount int) error {
					return errors.New("internal error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Unauthorized",
			request:        httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer([]byte(`{"toUser": "recipient", "amount": 100}`))),
			mockService:    &transaction.ServiceMock{},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.HandleSendCoin(tt.mockService)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, tt.request)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status code %v, got %v", tt.expectedStatus, status)
			}
		})
	}
}
