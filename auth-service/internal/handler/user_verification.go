package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/common/json"
)

func (h *Handler) UserVerification(w http.ResponseWriter, r *http.Request) {
	// read token from URL
	token := chi.URLParam(r, "token")
	if token == "" {
		json.WriteError(w, r, http.StatusBadRequest, "unable to parse token")
		return
	}

	// verify user
	if err := h.userService.VerifyUser(r.Context(), token); err != nil {
		json.WriteError(w, r, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
