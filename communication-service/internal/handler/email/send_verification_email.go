package email

import (
	"errors"
	"log"
	"net/http"

	"github.com/orduro/pos-microservices/common/json"
	"github.com/orduro/pos-microservices/communication-service/internal/store"
)

type VerificationEmailRequest struct {
	Email           string `json:"email"`
	VerificationURL string `json:"verification_url"`
}

func (h *Handler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	req := VerificationEmailRequest{}

	// read post request
	err := json.Read(r, &req)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	// send verification email
	err = h.store.Email.SendVerificationEmail(req.Email, req.VerificationURL)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrEmailInvalid):
			json.WriteError(w, r, http.StatusBadRequest, "invalid email format")
			return

		case errors.Is(err, store.ErrVerificationURLRequired):
			json.WriteError(w, r, http.StatusBadRequest, "verification URL is required")
			return

		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to send verification email")
			log.Printf("failed to send verification email: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
