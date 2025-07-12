package handler

import (
	"github.com/orduro/pos-microservices/auth-service/internal/httpclient"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type Handler struct {
	env              string
	frontendAdminURL string
	store            store.Store
	httpClient       *httpclient.Client
}

func New(env, frontendAdminURL string, s store.Store, httpClient *httpclient.Client) *Handler {
	return &Handler{
		env:              env,
		frontendAdminURL: frontendAdminURL,
		store:            s,
		httpClient:       httpClient,
	}
}
