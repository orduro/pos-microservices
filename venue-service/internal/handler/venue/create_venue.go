package venue

import (
	"errors"
	"log"
	"net/http"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/venue/internal/store"
)

func (h *Handler) CreateVenue(w http.ResponseWriter, r *http.Request) {
	venue := store.Venue{}

	// read post request
	err := json.Read(r, &venue)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	// validate venue
	if err := store.ValidateVenue(&venue); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "failed to validate venue")
		log.Printf("failed to validate venue: %v", err)
		return
	}

	// create venue
	if err := h.store.Venues.Create(r.Context(), &venue); err != nil {
		switch {
		case errors.Is(err, store.ErrVenueDuplicate):
			json.WriteError(w, r, http.StatusConflict, "venue with this name and address already exists")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to create venue")
			log.Printf("failed to create venue: %v", err)
			return
		}
	}

	json.Write(w, http.StatusCreated, venue)
}
