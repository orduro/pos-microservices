package handler

import "github.com/orduro/pos-microservices/auth-service/internal/store"

type Handler struct {
	env     string
	baseURL string
	store   store.Store
}

func New(env, baseURL string, s store.Store) *Handler {
	return &Handler{
		env:     env,
		baseURL: baseURL,
		store:   s,
	}
}
