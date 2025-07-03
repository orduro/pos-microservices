package venue

import "github.com/orduro/pos-microservices/venue/internal/store"

type Handler struct {
	store store.Store
}

func New(s store.Store) *Handler {
	return &Handler{
		store: s,
	}
}
