package venue

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/pos-microservices/common/json"
	"github.com/orduro/pos-microservices/venue/internal/store"
)

func (h *Handler) ArchiveVenue(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid venue id")
		return
	}

	err = h.store.Venues.Archive(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrVenueNotFound):
			json.WriteError(w, r, http.StatusNotFound, "venue not found")
			return
		case errors.Is(err, store.ErrVenueAlreadyArchived):
			json.WriteError(w, r, http.StatusConflict, "venue already archived")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "unable to archive venue")
			log.Printf("unable to archive venue: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteVenue(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid venue id")
		return
	}

	err = h.store.Venues.Delete(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrVenueNotFound):
			json.WriteError(w, r, http.StatusNotFound, "venue not found")
			return
		case errors.Is(err, store.ErrVenueNotArchived):
			json.WriteError(w, r, http.StatusBadRequest, "venue must be archived before deletion")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "unable to delete venue")
			log.Printf("unable to delete venue: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
