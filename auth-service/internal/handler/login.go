package handler

import (
	"log"
	"net/http"

	"github.com/orduro/pos-microservices/auth-service/internal/store"
	"github.com/orduro/pos-microservices/common/json"
)

type LoginResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDetails store.UserLoginDetails
	if err := json.Read(r, &loginDetails); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json")
		log.Printf("unable to read login details: %v", err)
		return
	}

	if err := loginDetails.Validate(); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid login details")
		log.Printf("invalid login details: %v", err)
		return
	}

	tokenPair, err := h.userService.LoginUser(r.Context(), loginDetails)
	if err != nil {
		switch err.Error() {
		case "invalid credentials":
			json.WriteError(w, r, http.StatusUnauthorized, "invalid email or password")
			return
		case "email not verified":
			json.WriteError(w, r, http.StatusForbidden, "please verify your email before logging in")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to login")
			log.Printf("failed to login user: %v", err)
			return
		}
	}

	response := LoginResponse{
		Message:      "login successful",
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}

	json.Write(w, http.StatusOK, response)
}
