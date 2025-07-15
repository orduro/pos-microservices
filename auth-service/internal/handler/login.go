package handler

import "net/http"

// TODO: get JWT token
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {}

// TODO: validate JWT for other services
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {}

// TODO: get current user info
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {}
