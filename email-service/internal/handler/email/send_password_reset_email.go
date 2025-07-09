package email

import (
	"errors"
	"log"
	"net/http"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/email-service/internal/store"
)

type PasswordResetEmailRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	ResetURL string `json:"reset_url"`
}

func (h *Handler) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	req := PasswordResetEmailRequest{}

	// read post request
	err := json.Read(r, &req)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	// send password reset email
	err = h.store.Email.SendPasswordResetEmail(req.Email, req.Username, req.ResetURL)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrEmailInvalid):
			json.WriteError(w, r, http.StatusBadRequest, "invalid email format")
			return

		case errors.Is(err, store.ErrUsernameRequired):
			json.WriteError(w, r, http.StatusBadRequest, "username is required")
			return

		case errors.Is(err, store.ErrResetURLRequired):
			json.WriteError(w, r, http.StatusBadRequest, "reset URL is required")
			return

		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to send password reset email")
			log.Printf("failed to send password reset email: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
