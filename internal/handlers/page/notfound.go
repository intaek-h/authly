package page

import (
	"net/http"

	"github.com/authly/internal/templates/pages"
)

type NotFoundPageHandler struct{}

func NewNotFoundPageHandler() *NotFoundPageHandler {
	return &NotFoundPageHandler{}
}

func (h *NotFoundPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := pages.NotFound().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "템플릿 제작중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}
}
