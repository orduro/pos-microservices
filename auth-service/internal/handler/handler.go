package handler

import (
	"github.com/orduro/pos-microservices/auth-service/internal/service"
)

type Handler struct {
	env         string
	userService *service.UserService
}

func New(env string, userService *service.UserService) *Handler {
	return &Handler{
		env:         env,
		userService: userService,
	}
}
