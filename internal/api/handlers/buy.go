package handlers

import (
	"errors"
	"github.com/nglmq/avito-shop/internal/app/merch"
	"net/http"
	"strings"
)

func HandleBuyItem(s merch.ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("user").(string)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "HandleBuyItem", ErrUnauthorized)
			return
		}

		item := strings.TrimPrefix(r.URL.Path, "/api/buy/")
		if item == "" {
			respondWithError(w, http.StatusBadRequest, "HandleBuyItem", ErrInvalidBody)
			return
		}

		err := s.BuyItem(r.Context(), username, item, 1)
		if err != nil {
			if errors.Is(err, merch.ErrItemNotFound) {
				respondWithError(w, http.StatusBadRequest, "HandleBuyItem", err)
				return
			}
			if errors.Is(err, merch.ErrInvalidAmount) {
				respondWithError(w, http.StatusBadRequest, "HandleBuyItem", err)
				return
			}

			respondWithError(w, http.StatusInternalServerError, "HandleBuyItem", err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
