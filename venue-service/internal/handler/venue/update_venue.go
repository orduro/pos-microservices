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

	// get existing venue from db
	existingVenue, err := h.store.Venues.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrVenueNotFound):
			json.WriteError(w, r, http.StatusNotFound, "venue not found")
			return

		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to retrieve venue")
			log.Printf("failed to retrieve venue with id %d: %v", id, err)
			return
		}
	}

	// check if venue is archived
	if existingVenue.Archived {
		json.WriteError(w, r, http.StatusConflict, "cannot update archived venue")
		return
	}

	// parse incoming update request (partial update)
	var updateRequest struct {
		Name        *string `json:"name,omitempty"`
		Address     *string `json:"address,omitempty"`
		Phone       *string `json:"phone,omitempty"`
		VenueType   *string `json:"venue_type,omitempty"`
		Description *string `json:"description,omitempty"`
	}

	err = json.Read(r, &updateRequest)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	// create updated venue by copying existing venue and applying changes
	updatedVenue := *existingVenue
	changes := []string{}

	// apply only the fields that were provided in the request
	if updateRequest.Name != nil && *updateRequest.Name != existingVenue.Name {
		updatedVenue.Name = *updateRequest.Name
		changes = append(changes, "name")
	}

	if updateRequest.Address != nil && *updateRequest.Address != existingVenue.Address {
		updatedVenue.Address = *updateRequest.Address
		changes = append(changes, "address")
	}

	if updateRequest.Phone != nil && *updateRequest.Phone != existingVenue.Phone {
		updatedVenue.Phone = *updateRequest.Phone
		changes = append(changes, "phone")
	}

	if updateRequest.VenueType != nil && *updateRequest.VenueType != existingVenue.VenueType {
		updatedVenue.VenueType = *updateRequest.VenueType
		changes = append(changes, "venue_type")
	}

	if updateRequest.Description != nil && *updateRequest.Description != existingVenue.Description {
		updatedVenue.Description = *updateRequest.Description
		changes = append(changes, "description")
	}

	// if no changes were made, return the existing venue
	if len(changes) == 0 {
		log.Printf("no changes detected for venue %d, returning existing venue", id)
		json.Write(w, http.StatusOK, existingVenue)
		return
	}

	log.Printf("updating venue %d with changes: %v", id, changes)

	// validate the updated venue
	if err := store.ValidateVenue(&updatedVenue); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "failed to validate venue")
		log.Printf("failed to validate venue: %v", err)
		return
	}

	// update venue in db
	if err := h.store.Venues.Update(r.Context(), &updatedVenue); err != nil {
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

	log.Printf("successfully updated venue %d with changes: %v", id, changes)
	json.Write(w, http.StatusOK, &updatedVenue)
}
