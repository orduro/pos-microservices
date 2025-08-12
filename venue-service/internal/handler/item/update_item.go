package item

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/pos-microservices/common/json"
	"github.com/orduro/pos-microservices/venue/internal/store"
)

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid item id")
		return
	}

	existingItem, err := h.store.Items.GetByID(r.Context(), int64(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrItemNotFound):
			json.WriteError(w, r, http.StatusNotFound, "item not found")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to retrieve item")
			log.Printf("failed to retrieve item with id %d: %v", id, err)
			return
		}
	}

	if existingItem.Archived {
		json.WriteError(w, r, http.StatusConflict, "cannot update archived item")
		return
	}

	var updateRequest struct {
		VenueID     *int64            `json:"venue_id"`
		Name        *string           `json:"name"`
		Description *string           `json:"description"`
		Category    *string           `json:"category"`
		Price       *float64          `json:"price"`
		IsAvailable *bool             `json:"is_available"`
		ImageURL    *string           `json:"image_url"`
		Position    *int              `json:"position"`
		Modifiers   *[]store.Modifier `json:"modifiers"`
	}

	err = json.Read(r, &updateRequest)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	updatedItem := *existingItem
	changes := []string{}

	if updateRequest.VenueID != nil && *updateRequest.VenueID != existingItem.VenueID {
		updatedItem.VenueID = *updateRequest.VenueID
		changes = append(changes, "venue_id")
	}

	if updateRequest.Name != nil && *updateRequest.Name != existingItem.Name {
		updatedItem.Name = *updateRequest.Name
		changes = append(changes, "name")
	}

	if updateRequest.Description != nil && *updateRequest.Description != existingItem.Description {
		updatedItem.Description = *updateRequest.Description
		changes = append(changes, "description")
	}

	if updateRequest.Category != nil && *updateRequest.Category != existingItem.Category {
		updatedItem.Category = *updateRequest.Category
		changes = append(changes, "category")
	}

	if updateRequest.Price != nil && *updateRequest.Price != existingItem.Price {
		updatedItem.Price = *updateRequest.Price
		changes = append(changes, "price")
	}

	if updateRequest.IsAvailable != nil && *updateRequest.IsAvailable != existingItem.IsAvailable {
		updatedItem.IsAvailable = *updateRequest.IsAvailable
		changes = append(changes, "is_available")
	}

	if updateRequest.ImageURL != nil && *updateRequest.ImageURL != existingItem.ImageURL {
		updatedItem.ImageURL = *updateRequest.ImageURL
		changes = append(changes, "image_url")
	}

	if updateRequest.Position != nil && *updateRequest.Position != existingItem.Position {
		updatedItem.Position = *updateRequest.Position
		changes = append(changes, "position")
	}

	if updateRequest.Modifiers != nil {
		updatedItem.Modifiers = *updateRequest.Modifiers
		changes = append(changes, "modifiers")
	}

	if len(changes) == 0 {
		log.Printf("no changes detected for item %d, returning existing item", id)
		json.Write(w, http.StatusOK, existingItem)
		return
	}

	log.Printf("updating item %d with changes: %v", id, changes)

	if err := store.ValidateItem(&updatedItem); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "failed to validate item")
		log.Printf("failed to validate item: %v", err)
		return
	}

	if err := h.store.Items.Update(r.Context(), &updatedItem); err != nil {
		switch {
		case errors.Is(err, store.ErrItemNotFound):
			json.WriteError(w, r, http.StatusNotFound, "item not found")
			return
		case errors.Is(err, store.ErrItemAlreadyArchived):
			json.WriteError(w, r, http.StatusConflict, "cannot update archived item")
			return
		case errors.Is(err, store.ErrItemDuplicate):
			json.WriteError(w, r, http.StatusConflict, "item with this name already exists in venue")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to update item")
			log.Printf("failed to update item: %v", err)
			return
		}
	}

	log.Printf("successfully updated item %d with changes: %v", id, changes)
	json.Write(w, http.StatusOK, &updatedItem)
}
