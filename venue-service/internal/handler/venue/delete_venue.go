package venue

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/common/json"
)

func (h *Handler) ArchiveVenue(w http.ResponseWriter, r *http.Request) {
	// parse id from URL
	urlID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, fmt.Sprintf("cannot parse venue id: %v", urlID))
		log.Printf("cannot parse venue id %v", urlID)
		return
	}

	// update the venue with the corresponding id
	if err := h.store.Venues.Archive(r.Context(), id); err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, "unable to archive venue")
		log.Printf("unable to archive venue: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
