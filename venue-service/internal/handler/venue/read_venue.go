package venue

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/common/json"
)

func (h *Handler) GetVenueByID(w http.ResponseWriter, r *http.Request) {
	// read id from URL
	urlID := chi.URLParam(r, "id")

	// convert to correct type
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, fmt.Sprintf("cannot parse venue id: %v", urlID))
		return
	}

	// get venue with the id requested
	venue, err := h.store.Venues.GetByID(r.Context(), id)
	if err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, fmt.Sprintf("cannot find venue with id: %v", id))
	}

	// write response
	json.Write(w, http.StatusOK, venue)
}
