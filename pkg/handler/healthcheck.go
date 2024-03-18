package handler

import (
	"net/http"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	h.renderer.RenderOK(w)
}
