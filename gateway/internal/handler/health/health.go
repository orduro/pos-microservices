package health

import (
	"net/http"
	"time"

	"github.com/orduro/common/json"
)

type Health struct {
	Service     string    `json:"service"`
	Health      string    `json:"health"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"timestamp"`
}

// @Summary Health check endpoint
// @Description Get the health status of the gateway service
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Health
// @Router /api/health [get]
func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	data := &Health{
		Service:     "gateway",
		Health:      "alive",
		Environment: h.env,
		Timestamp:   time.Now().Local(),
	}
	json.Write(w, http.StatusOK, data)
}
