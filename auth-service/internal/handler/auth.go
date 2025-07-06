package handler

import "net/http"

// create user
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {}

// verify user (send email)
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {}

// get JWT token
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {}

// validate JWT for other services
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {}

// get current user info
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {}
