package handler

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/renderer"
)

type Handler struct {
	db       *db.Database
	renderer renderer.Renderer

	url string
}

func NewHandler(db *db.Database, url string) *Handler {
	return &Handler{
		db:  db,
		url: fmt.Sprintf("%s/api/v1", url),
	}
}

func (h *Handler) Router() chi.Router {
	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {

		r.Route("/admin", func(r chi.Router) {
			r.Delete("/{secretKey}", h.Delete)
			r.Get("/{secretKey}", h.Info)
		})

		r.Post("/make_shorter", h.Short)
		r.Get("/{shortSuffix}", h.Redirect)

		r.Get("/healthcheck/ping", h.HealthCheck)
	})

	return router
}
