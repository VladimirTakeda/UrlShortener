package renderer

import (
	"encoding/json"
	"errors"
	"net/http"

	apierrors "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/errors"
)

type Renderer struct{}

func (r Renderer) RenderJSON(w http.ResponseWriter, response interface{}) {
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(response)
}

func (r Renderer) RenderOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (r Renderer) RenderError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	var errorWithCode apierrors.Error
	if errors.As(err, &errorWithCode) {
		status = errorWithCode.Code()
	}

	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	_ = encoder.Encode(ErrorResponse{
		Error: err.Error(),
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}
