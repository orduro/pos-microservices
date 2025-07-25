package venue

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/venue/internal/store"
)

type PaginationMetadata struct {
	// the current page number
	CurrentPage int `json:"current_page"`

	// number of records requested per page (limit)
	PerPage int `json:"per_page"`

	// total number of pages available
	TotalPages int `json:"total_pages"`

	// total number of venues matching the filter criteria
	TotalRecords int `json:"total_records"`

	// indicates if there are more pages after the current one
	HasNext bool `json:"has_next"`

	// indicates if there are page before current one
	HasPrevious bool `json:"has_previous"`

	// actual number of venues returned in the response
	Showing int `json:"showing"`
}

type VenueListResponse struct {
	Venues   []store.Venue      `json:"venues"`
	Metadata PaginationMetadata `json:"metadata"`
}

// ListVenues retrieves a paginated list of venues for the authenticated tenant
//
// Query Parameters:
//
//   - include_archived (optional): Include archived venues in results (default: false)
//     Values: "true" | "false"
//     Example: ?include_archived=true
//
//   - search (optional): Search term to filter venues by name, address, or description
//     Type: string
//     Example: ?search=restaurant
//
//   - limit (optional): Number of venues per page (default: 50, max: 100)
//     Type: integer (1-100)
//     Example: ?limit=25
//
//   - offset (optional): Number of venues to skip for pagination (default: 0)
//     Type: integer (>= 0)
//     Example: ?offset=50
//
// Response Format:
//   - venues: Array of venue objects
//   - metadata: Pagination information including current_page, total_pages, etc.
//
// Examples:
//
//	GET /api/venue                                   - Get first 50 active venues
//	GET /api/venue?limit=10&offset=20                - Get venues 21-30
//	GET /api/venue?search=cafe&include_archived=true - Search for "cafe" including archived
//	GET /api/venue?limit=25&search=restaurant        - Search for "restaurant" with 25 per page
func (h *Handler) ListVenues(w http.ResponseWriter, r *http.Request) {
	// parse query parameters for filtering and pagination
	query := r.URL.Query()

	// get include_archived parameter (default: false)
	includeArchived := query.Get("include_archived") == "true"

	// get pagination parameters
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

	// get search parameter
	search := query.Get("search")

	// create filter options
	filter := store.VenueFilter{
		IncludeArchived: includeArchived,
		Search:          search,
		Limit:           limit,
		Offset:          offset,
	}

	// get venues from store
	venues, total, err := h.store.Venues.List(r.Context(), filter)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrTenantNotFound):
			json.WriteError(w, r, http.StatusBadRequest, "tenant not found")
			return
		default:
			json.WriteError(w, r, http.StatusInternalServerError, "failed to retrieve venues")
			log.Printf("failed to retrieve venues: %v", err)
			return
		}
	}

	// calculate pagination metadata
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	currentPage := (offset / limit) + 1
	showing := len(venues)
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

	response := VenueListResponse{
		Venues:   venues,
		Metadata: metadata,
	}

	json.Write(w, http.StatusOK, response)
}
