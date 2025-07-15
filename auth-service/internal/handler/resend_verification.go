package handler

import (
	"log"
	"net/http"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (r *ResendVerificationRequest) Validate() error {
	return store.V.Struct(r)
}

func (h *Handler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
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

	err := h.userService.ResendVerificationEmail(r.Context(), req.Email)
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
