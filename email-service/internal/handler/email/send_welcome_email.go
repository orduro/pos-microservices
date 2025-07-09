package email

import (
	"errors"
	"log"
	"net/http"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/email-service/internal/store"
)

type WelcomeEmailRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (h *Handler) SendWelcomeEmail(w http.ResponseWriter, r *http.Request) {
	req := WelcomeEmailRequest{}

	// read post request
	err := json.Read(r, &req)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	// send welcome email
	err = h.store.Email.SendWelcomeEmail(req.Email, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrEmailInvalid):
			json.WriteError(w, r, http.StatusBadRequest, "invalid email format")
			return

		case errors.Is(err, store.ErrUsernameRequired):
			json.WriteError(w, r, http.StatusBadRequest, "username is required")
			return

		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to send welcome email")
			log.Printf("failed to send welcome email: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
