package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	database "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
	apierrors "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/errors"
)

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shortSuffix := chi.URLParam(r, "shortSuffix")

	fmt.Printf(shortSuffix + "\n")

	link, err := h.db.SelectBySuffix(ctx, shortSuffix)
	switch {
	case errors.Is(err, database.SuffixNotFoundError):
		fmt.Printf("failed to found short suffix \"%s\"\n", shortSuffix)
		h.renderer.RenderError(w, apierrors.NotFoundError{})
		return
	case err != nil:
		fmt.Printf("failed to select link by short suffix \"%s\"\n", shortSuffix)
		h.renderer.RenderError(w, apierrors.NotFoundError{})
		return
	}

	err = h.db.IncrementClicksBySuffix(ctx, shortSuffix)
	if err != nil {
		fmt.Printf("failed to increment clicks by short suffix \"%s\"\n", shortSuffix)
		h.renderer.RenderError(w, apierrors.NotFoundError{})
		return
	}

	http.Redirect(w, r, link.Link, http.StatusTemporaryRedirect)
}
