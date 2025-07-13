package handler

import (
	"net/http"
	"time"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/auth-service/internal/constants"
)

type Health struct {
	Service     string    `json:"service"`
	Health      string    `json:"health"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"timestamp"`
}

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	data := &Health{
		Service:     constants.ServiceNameAuth,
		Health:      constants.StatusHealthy,
		Environment: h.env,
		Timestamp:   time.Now().Local(),
	}
	json.Write(w, http.StatusOK, data)
}
