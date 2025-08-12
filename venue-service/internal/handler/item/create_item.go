package item

import (
	"errors"
	"log"
	"net/http"

	"github.com/orduro/pos-microservices/common/json"
	"github.com/orduro/pos-microservices/venue/internal/store"
)

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	item := store.Item{}

	err := json.Read(r, &item)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid json value, cannot be read")
		log.Printf("invalid json value, cannot be read: %v", err)
		return
	}

	if err := store.ValidateItem(&item); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "failed to validate item")
		log.Printf("failed to validate item: %v", err)
		return
	}

	if err := h.store.Items.Create(r.Context(), &item); err != nil {
		switch {
		case errors.Is(err, store.ErrItemDuplicate):
			json.WriteError(w, r, http.StatusConflict, "item with this name already exists in venue")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to create item")
			log.Printf("failed to create item: %v", err)
			return
		}
	}

	json.Write(w, http.StatusCreated, item)
}
