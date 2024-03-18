package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	database "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
	apierrors "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/errors"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	secretKey := chi.URLParam(r, "secretKey")

	err := h.db.DeleteBySecretKey(ctx, secretKey)
	switch {
	case errors.Is(err, database.SecretKeyNotFoundError):
		fmt.Printf("failed to found secret key \"%s\"\n", secretKey)
		h.renderer.RenderError(w, apierrors.NotFoundError{})
		return
	case err != nil:
		fmt.Printf("failed to delete secret key\n")
		h.renderer.RenderError(w, apierrors.NotFoundError{})
		return
	}

	h.renderer.RenderOK(w)
}
