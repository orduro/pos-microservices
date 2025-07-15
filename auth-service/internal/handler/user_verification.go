package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

func (h *Handler) UserVerification(w http.ResponseWriter, r *http.Request) {
	// read token from URL
	token := chi.URLParam(r, "token")
	if token == "" {
		json.WriteError(w, r, http.StatusBadRequest, "unable to parse token")
		return
	}

	// validate token format
	if len(token) < constants.TokenByteLength {
		json.WriteError(w, r, http.StatusBadRequest, "invalid verification token format")
		return
	}

	// verify user
	err := h.userService.VerifyUser(r.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrTokenNotFound):
			json.WriteError(w, r, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, store.ErrTokenExpired):
			json.WriteError(w, r, http.StatusGone, err.Error())
			return
		case err.Error() == "user is already verified":
			json.WriteError(w, r, http.StatusConflict, "user is already verified")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "something unexpected happened")
			log.Printf("unexpected error during user verification: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
