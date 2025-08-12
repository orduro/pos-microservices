package item

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/venue/internal/store"
)

func (h *Handler) ArchiveItem(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid item id")
		return
	}

	err = h.store.Items.Archive(r.Context(), int64(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrItemNotFound):
			json.WriteError(w, r, http.StatusNotFound, "item not found")
			return
		case errors.Is(err, store.ErrItemAlreadyArchived):
			json.WriteError(w, r, http.StatusConflict, "item already archived")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "unable to archive item")
			log.Printf("unable to archive item: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid item id")
		return
	}

	err = h.store.Items.Delete(r.Context(), int64(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrItemNotFound):
			json.WriteError(w, r, http.StatusNotFound, "item not found")
			return
		case errors.Is(err, store.ErrItemNotArchived):
			json.WriteError(w, r, http.StatusBadRequest, "item must be archived before deletion")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "unable to delete item")
			log.Printf("unable to delete item: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
