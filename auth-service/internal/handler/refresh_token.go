package handler

import (
	"net/http"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (r *RefreshTokenRequest) Validate() error {
	return store.V.Struct(r)
}

type RefreshTokenResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json")
		return
	}

	if err := req.Validate(); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "refresh token is required")
		return
	}

	tokenPair, err := h.userService.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		json.WriteError(w, r, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	// don't include new refresh token
	// the only time client gets refreshtoken is
	// during sign in
	response := RefreshTokenResponse{
		Message:      "tokens refreshed successfully",
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: req.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
	}

	json.Write(w, http.StatusOK, response)
}
