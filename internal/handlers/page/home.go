package page

import (
	"context"
	"net/http"

	"github.com/authly/internal/templates/pages"
)

func (p *Pages) HandlerHomePage(w http.ResponseWriter, r *http.Request) {
	page := pages.Home()

	page.Render(createHomePageContext(r.Context()), w)

	// page := pages.Home()

	// if r.Header.Get("HX-Request") == "true" {
	// 	err := page.Render(r.Context(), w)
	// 	if err != nil {
	// 		http.Error(w, "템플릿 제작중 오류가 발생했습니다.", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	return
	// }

	// session, err := p.Session.Get(r, "auth-session")

	// if session.IsNew || err != nil {
	// 	err := layouts.DefaultLayout(page, "홈페이지", env.MustLoad()).Render(r.Context(), w)
	// 	if err != nil {
	// 		http.Error(w, "템플릿 제작중 오류가 발생했습니다.", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	return
	// }

	// profile := session.Values["profile"].(map[string]interface{})

	// err = layouts.DefaultLayout(page, profile["nickname"].(string), env.MustLoad()).Render(r.Context(), w)

	// if err != nil {
	// 	http.Error(w, "템플릿 제작중 오류가 발생했습니다.", http.StatusInternalServerError)
	// 	return
	// }
}

func createHomePageContext(ctx context.Context) context.Context {
	pageTitle := "홈페이지"

	c := context.WithValue(ctx, "head-title", pageTitle)

	return c
}
