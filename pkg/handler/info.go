package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	database "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
	apierrors "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/errors"
)

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	secretKey := chi.URLParam(r, "secretKey")

	fmt.Printf("got secret key \"%s\"\n", secretKey)

	link, err := h.db.SelectBySecretKey(ctx, secretKey)
	switch {
	case errors.Is(err, database.SecretKeyNotFoundError):
		fmt.Printf("failed to found secret key \"%s\"\n", secretKey)
		h.renderer.RenderError(w, apierrors.NotFoundError{})
		return
	case err != nil:
		fmt.Printf("failed to select secret key\n")
		h.renderer.RenderError(w, apierrors.NotFoundError{})
		return
	}

	var response InfoResponse
	response.FromLink(link, h.url)

	h.renderer.RenderJSON(w, response)
}
