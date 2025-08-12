package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
	"github.com/orduro/pos-microservices/common/json"
)

func (h *Handler) UserVerificationCallback(w http.ResponseWriter, r *http.Request) {
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

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (r *ResendVerificationRequest) Validate() error {
	return store.V.Struct(r)
}

func (h *Handler) SendEmailVerification(w http.ResponseWriter, r *http.Request) {
	var req ResendVerificationRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json")
		log.Printf("unable to read resend verification request: %v", err)
		return
	}

	if err := req.Validate(); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid email")
		log.Printf("invalid email format: %v", err)
		return
	}

	err := h.userService.SendEmailVerification(r.Context(), req.Email)
	if err != nil {
		if err.Error() == "user is already verified" {
			json.WriteError(w, r, http.StatusConflict, err.Error())
			return
		}
		json.WriteError(w, r, http.StatusInternalServerError, "failed to resend verification email")
		log.Printf("failed to resend verification email: %v", err)
		return
	}

	json.Write(w, http.StatusOK, map[string]any{
		"message": constants.MsgVerificationSent,
	})
}
