package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/nglmq/avito-shop/internal/app/transaction"
	"github.com/nglmq/avito-shop/internal/models"
	"net/http"
)

func HandleSendCoin(s transaction.ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("user").(string)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "HandleSendCoin", ErrUnauthorized)
			return
		}

		var req models.SendCoinsRequest
		validate := validator.New()

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "HandleSendCoin", ErrInvalidBody)
			return
		}
		err := validate.Struct(req)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "HandleAuth", ErrInvalidBody)
			return
		}

		errCh := make(chan error)

		go func() {
			err := s.SendCoins(r.Context(), username, req.ToUser, req.Amount)
			errCh <- err
		}()

		err = <-errCh
		if err != nil {
			if errors.Is(err, transaction.ErrInvalidAmount) {
				respondWithError(w, http.StatusBadRequest, "HandleSendCoin", err)
				return
			}
			if errors.Is(err, transaction.ErrInvalidRecipient) {
				respondWithError(w, http.StatusBadRequest, "HandleSendCoin", err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, "HandleSendCoin", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
