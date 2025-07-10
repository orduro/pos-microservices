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

func (h *Handler) UpdateVenue(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid venue id")
		return
	}

	venue := store.Venue{ID: id}

	// read put request
	err = json.Read(r, &venue)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	venue.ID = id

	// validate venue
	if err := store.ValidateVenue(&venue); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "failed to validate venue")
		log.Printf("failed to validate venue: %v", err)
		return
	}

	// update venue
	if err := h.store.Venues.Update(r.Context(), &venue); err != nil {
		switch {
		case errors.Is(err, store.ErrVenueNotFound):
			json.WriteError(w, r, http.StatusNotFound, "venue not found")
			return

		case errors.Is(err, store.ErrVenueAlreadyArchived):
			json.WriteError(w, r, http.StatusConflict, "cannot update archived venue")
			return

		case errors.Is(err, store.ErrVenueDuplicate):
			json.WriteError(w, r, http.StatusConflict, "venue with this name and address already exists")
			return

		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to update venue")
			log.Printf("failed to update venue: %v", err)
			return
		}
	}

	json.Write(w, http.StatusOK, venue)
}
