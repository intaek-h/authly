package page

import (
	"net/http"

	"github.com/authly/internal/env"
	"github.com/authly/internal/templates/layouts"
	"github.com/authly/internal/templates/pages"
)

type HomePageHandler struct{}

func NewHomePageHandler() *HomePageHandler {
	return &HomePageHandler{}
}

func (h *HomePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := pages.Home()
	err := layouts.DefaultLayout(page, "인택", env.MustLoad()).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "템플릿 제작중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}
}