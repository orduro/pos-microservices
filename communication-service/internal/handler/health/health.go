package health

import (
	"net/http"
	"time"

	"github.com/orduro/pos-microservices/common/json"
)

type Health struct {
	Service     string    `json:"service"`
	Health      string    `json:"health"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"timestamp"`
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	data := &Health{
		Service:     "communication-service",
		Health:      "alive",
		Environment: h.env,
		Timestamp:   time.Now().Local(),
	}
	json.Write(w, http.StatusOK, data)
}
