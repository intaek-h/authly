package page

import (
	"context"
	"net/http"

	"github.com/authly/internal/templates/pages"
)

type NotFoundPageHandler struct{}

func NewNotFoundPageHandler() *NotFoundPageHandler {
	return &NotFoundPageHandler{}
}

func (p *Pages) HandlerNotFoundPage(w http.ResponseWriter, r *http.Request) {
	page := pages.NotFound()
	page.Render(createNotFoundPageContext(r.Context()), w)

	// if r.Header.Get("HX-Request") == "true" {
	// 	err := page.Render(r.Context(), w)
	// 	if err != nil {
	// 		http.Error(w, "템플릿 제작중 오류가 발생했습니다.", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	return
	// }

	// err := layouts.DefaultLayout(page, "인택", env.MustLoad()).Render(r.Context(), w)
	// if err != nil {
	// 	http.Error(w, "템플릿 제작중 오류가 발생했습니다.", http.StatusInternalServerError)
	// 	return
	// }
}

func createNotFoundPageContext(ctx context.Context) context.Context {
	pageTitle := "페이지를 찾을 수 없습니다."

	c := context.WithValue(ctx, "head-title", pageTitle)

	return c
}
