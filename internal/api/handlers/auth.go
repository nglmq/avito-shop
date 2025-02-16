package handlers

import (
	"encoding/json"
	"github.com/nglmq/avito-shop/internal/app/auth"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nglmq/avito-shop/internal/models"
)

func HandleAuth(service auth.ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AuthRequest
		validate := validator.New()

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "HandleAuth", ErrInvalidBody)
			return
		}

		err = validate.Struct(req)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "HandleAuth", ErrInvalidBody)
			return
		}

		if req.Username == "" || req.Password == "" {
			respondWithError(w, http.StatusBadRequest, "HandleAuth", ErrInvalidBody)
			return
		}

		token, err := service.AuthenticateUser(r.Context(), req.Username, req.Password)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "HandleAuth", ErrUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(models.AuthResponse{Token: token}); err != nil {
			return
		}
	}
}
