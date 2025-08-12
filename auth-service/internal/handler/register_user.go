package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
	"github.com/orduro/pos-microservices/common/json"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registrationDetails store.UserRegistrationDetails
	if err := json.Read(r, &registrationDetails); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json")
		log.Printf("unable to read registration details: %v", err)
		return
	}

	if err := registrationDetails.Validate(); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid registration details")
		log.Printf("invalid registration details: %v", err)
		return
	}

	// create user in system
	user, err := h.userService.RegisterUser(r.Context(), registrationDetails)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrUserExists):
			json.WriteError(w, r, http.StatusConflict, err.Error())
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to register user")
			log.Printf("failed to register user: %v", err)
			return
		}
	}

	// send verification email
	err = h.userService.SendEmailVerification(r.Context(), user.Email)
	if err != nil {
		if err.Error() == "user is already verified" {
			json.WriteError(w, r, http.StatusConflict, err.Error())
			return
		}
		json.WriteError(w, r, http.StatusInternalServerError, "failed to resend verification email")
		log.Printf("failed to resend verification email: %v", err)
		return
	}

	json.Write(w, http.StatusCreated, map[string]any{
		"message": constants.MsgRegistrationSuccess,
		"email":   user.Email,
	})
}
