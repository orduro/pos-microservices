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

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (f *ForgotPasswordRequest) Validate() error {
	return store.V.Struct(f)
}

func (h *Handler) InitiatePasswordReset(w http.ResponseWriter, r *http.Request) {
	// get email from request
	var req ForgotPasswordRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json")
		return
	}

	// validate req
	if err := req.Validate(); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid email")
		return
	}

	if err := h.userService.SendForgetPasswordEmail(r.Context(), req.Email); err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

type ForgotPasswordCallbackRequest struct {
	Password string `json:"password" validate:"required,min=8,max=128"`
}

func (f *ForgotPasswordCallbackRequest) Validate() error {
	return store.V.Struct(f)
}

func (h *Handler) PasswordResetCallback(w http.ResponseWriter, r *http.Request) {
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

	// get password from body
	var req ForgotPasswordCallbackRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json")
		return
	}

	// validate password
	if err := req.Validate(); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "password validation failed")
		return
	}

	// reset password
	err := h.userService.ResetPassword(r.Context(), token, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrTokenNotFound):
			json.WriteError(w, r, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, store.ErrTokenExpired):
			json.WriteError(w, r, http.StatusGone, err.Error())
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "something unexpected happened")
			log.Printf("unexpected error during user verification: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
