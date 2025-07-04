package venue

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/venue/internal/store"
)

func (h *Handler) GetVenueByID(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid venue id")
		return
	}

	venue, err := h.store.Venues.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrVenueNotFound) {
			json.WriteError(w, r, http.StatusNotFound, "venue not found")
			return
		}
		json.WriteError(w, r, http.StatusInternalServerError, "failed to retrieve venue")
		log.Printf("failed to retrieve venue with id %d: %v", id, err)
		return
	}

	json.Write(w, http.StatusOK, venue)
}
