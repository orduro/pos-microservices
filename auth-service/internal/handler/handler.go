package handler

import (
	"github.com/orduro/pos-microservices/auth-service/internal/httpclient"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type Handler struct {
	env        string
	baseURL    string
	store      store.Store
	httpClient *httpclient.Client
}

func New(env, baseURL string, s store.Store, httpClient *httpclient.Client) *Handler {
	return &Handler{
		env:        env,
		baseURL:    baseURL,
		store:      s,
		httpClient: httpClient,
	}
}

