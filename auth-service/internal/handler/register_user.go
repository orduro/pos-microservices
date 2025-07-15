package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
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

	json.Write(w, http.StatusCreated, map[string]any{
		"message": constants.MsgRegistrationSuccess,
		"email":   user.Email,
	})
}

// TODO: verify user (send email)
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {}

// TODO: get JWT token
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {}

// TODO: validate JWT for other services
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {}

// TODO: get current user info
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {}
