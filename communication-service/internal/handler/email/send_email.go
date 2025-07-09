package email

import (
	"errors"
	"log"
	"net/http"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/communication-service/internal/store"
)

func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	emailReq := store.EmailRequest{}

	// read post request
	err := json.Read(r, &emailReq)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	// send email
	response, err := h.store.Email.SendEmail(emailReq)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrEmailInvalid):
			json.WriteError(w, r, http.StatusBadRequest, "invalid email format")
			return

		case errors.Is(err, store.ErrSubjectRequired):
			json.WriteError(w, r, http.StatusBadRequest, "subject is required")
			return

		case errors.Is(err, store.ErrBodyRequired):
			json.WriteError(w, r, http.StatusBadRequest, "body is required")
			return

		case errors.Is(err, store.ErrRecipientsRequired):
			json.WriteError(w, r, http.StatusBadRequest, "at least one recipient is required")
			return

		case errors.Is(err, store.ErrEmailTypeInvalid):
			json.WriteError(w, r, http.StatusBadRequest, "invalid email type")
			return

		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to send email")
			log.Printf("failed to send email: %v", err)
			return
		}
	}

	json.Write(w, http.StatusOK, response)
}
