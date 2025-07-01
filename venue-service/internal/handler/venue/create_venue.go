package venue

import (
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
		json.WriteError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// create venue
	if err := h.store.Venues.Create(r.Context(), &venue); err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, "failed to create venue")
		log.Printf("failed to create venue: %v", err)
		return
	}

	json.Write(w, http.StatusCreated, venue)
}
