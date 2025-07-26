package handler

import (
	"net/http"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (f *ForgotPasswordRequest) Validate() error {
	return store.V.Struct(f)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
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
