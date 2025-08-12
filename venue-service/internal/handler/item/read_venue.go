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

type PaginationMetadata struct {
	CurrentPage  int  `json:"current_page"`
	PerPage      int  `json:"per_page"`
	TotalPages   int  `json:"total_pages"`
	TotalRecords int  `json:"total_records"`
	HasNext      bool `json:"has_next"`
	HasPrevious  bool `json:"has_previous"`
	Showing      int  `json:"showing"`
}

type ItemListResponse struct {
	Items    []store.Item       `json:"items"`
	Metadata PaginationMetadata `json:"metadata"`
}

// retrieves a single menu item
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := h.store.Items.GetByID(r.Context(), int64(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrItemNotFound):
			json.WriteError(w, r, http.StatusNotFound, "item not found")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to retrieve item")
			log.Printf("failed to retrieve item: %v", err)
			return
		}
	}

	json.Write(w, http.StatusOK, item)
}

// lists menu items with filtering and pagination
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	includeArchived := query.Get("include_archived") == "true"
	onlyAvailable := query.Get("only_available") == "true"

	var venueID int64
	if venueIDStr := query.Get("venue_id"); venueIDStr != "" {
		if parsedVenueID, err := strconv.ParseInt(venueIDStr, 10, 64); err == nil {
			venueID = int64(parsedVenueID)
		}
	}

	limit := 50
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := query.Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	category := query.Get("category")
	search := query.Get("search")

	filter := store.ItemFilter{
		VenueID:         venueID,
		IncludeArchived: includeArchived,
		OnlyAvailable:   onlyAvailable,
		Category:        category,
		Search:          search,
		Limit:           limit,
		Offset:          offset,
	}

	items, total, err := h.store.Items.List(r.Context(), filter)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrTenantNotFound):
			json.WriteError(w, r, http.StatusBadRequest, "tenant not found")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to retrieve items")
			log.Printf("failed to retrieve items: %v", err)
			return
		}
	}

	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	currentPage := (offset / limit) + 1
	showing := len(items)
	hasNext := offset+limit < total
	hasPrevious := offset > 0

	metadata := PaginationMetadata{
		CurrentPage:  currentPage,
		PerPage:      limit,
		TotalPages:   totalPages,
		TotalRecords: total,
		HasNext:      hasNext,
		HasPrevious:  hasPrevious,
		Showing:      showing,
	}

	response := ItemListResponse{
		Items:    items,
		Metadata: metadata,
	}

	json.Write(w, http.StatusOK, response)
}

// lists items for a specific venue
func (h *Handler) ListItemsByVenue(w http.ResponseWriter, r *http.Request) {
	urlVenueID := chi.URLParam(r, "venue_id")
	venueID, err := strconv.ParseInt(urlVenueID, 10, 64)
	if err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "invalid venue id")
		return
	}

	query := r.URL.Query()
	includeArchived := query.Get("include_archived") == "true"
	onlyAvailable := query.Get("only_available") != "false"

	limit := 50
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := query.Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	category := query.Get("category")
	search := query.Get("search")

	filter := store.ItemFilter{
		VenueID:         int64(venueID),
		IncludeArchived: includeArchived,
		OnlyAvailable:   onlyAvailable,
		Category:        category,
		Search:          search,
		Limit:           limit,
		Offset:          offset,
	}

	items, total, err := h.store.Items.ListByVenue(r.Context(), int64(venueID), filter)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrTenantNotFound):
			json.WriteError(w, r, http.StatusBadRequest, "tenant not found")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to retrieve items")
			log.Printf("failed to retrieve items: %v", err)
			return
		}
	}

	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	currentPage := (offset / limit) + 1
	showing := len(items)
	hasNext := offset+limit < total
	hasPrevious := offset > 0

	metadata := PaginationMetadata{
		CurrentPage:  currentPage,
		PerPage:      limit,
		TotalPages:   totalPages,
		TotalRecords: total,
		HasNext:      hasNext,
		HasPrevious:  hasPrevious,
		Showing:      showing,
	}

	response := ItemListResponse{
		Items:    items,
		Metadata: metadata,
	}

	json.Write(w, http.StatusOK, response)
}
