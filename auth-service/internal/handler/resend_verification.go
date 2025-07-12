package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (r *ResendVerificationRequest) Validate() error {
	return store.V.Struct(r)
}

func (h *Handler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	// parse request
	var req ResendVerificationRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	// validate request data
	if err := req.Validate(); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid email format")
		log.Printf("invalid email format: %v", err)
		return
	}

	// check if user exists
	user, err := h.store.Users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrUserNotFound):
			json.Write(w, http.StatusOK, map[string]any{
				"message": "if an account with this email exists, a verification email has been sent.",
			})
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "unable to process request")
			log.Printf("failed to check user: %v", err)
			return
		}
	}

	// check if user is already verified
	if user.IsVerified {
		json.WriteError(w, r, http.StatusConflict, "user is already verified")
		return
	}

	// generate new verification token
	token, err := generateVerificationToken()
	if err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, "unable to generate verification token")
		log.Printf("failed to generate verification token: %v", err)
		return
	}

	// create verification token
	verificationToken := store.VerificationToken{
		UserID:    user.ID,
		Token:     token,
		TokenType: "email_verification",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := h.store.VerificationTokens.Create(r.Context(), &verificationToken); err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, "unable to create verification token")
		log.Printf("failed to create verification token: %v", err)
		return
	}

	// send verification email in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		verificationLink := fmt.Sprintf("%s/verification?token=%s", h.frontendAdminURL, token)

		if err := h.httpClient.SendVerificationEmail(ctx, req.Email, verificationLink); err != nil {
			log.Printf("failed to send verification email to %s: %v", req.Email, err)
		} else {
			log.Printf("verification email sent successfully to %s", req.Email)
		}
	}()

	json.Write(w, http.StatusOK, map[string]any{
		"message": "if an account with this email exists, a verification email has been sent.",
	})
}

