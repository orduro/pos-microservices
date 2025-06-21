package health

import (
	"net/http"
	"time"

	jsoncodec "github.com/orduro/pos-microservices/pkg/json"
)

type Handler struct {
	env string
}

func New(env string) *Handler {
	return &Handler{
		env: env,
	}
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Health      string    `json:"health"`
		Environment string    `json:"environment"`
		Timestamp   time.Time `json:"timestamp"`
	}{
		Health:      "alive",
		Environment: h.env,
		Timestamp:   time.Now().Local(),
	}
	jsoncodec.Write(w, http.StatusOK, data)
}
