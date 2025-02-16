package handlers

import (
	"encoding/json"
	"github.com/nglmq/avito-shop/internal/app/history"
	"net/http"
)

func HandleGetInfo(s history.InfoServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("user").(string)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, "HandleGetInfo", ErrUnauthorized)
			return
		}

		info, err := s.GetInfo(r.Context(), username)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "HandleGetInfo", ErrInternal)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(info); err != nil {
			respondWithError(w, http.StatusInternalServerError, "HandleGetInfo", ErrInternal)
			return
		}
	}
}
