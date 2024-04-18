package page

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func (p *Pages) HandlerLoginPage(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		http.Error(w, "로그인 상태를 생성하는 중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "state",
		Value:    state,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	})

	// Redirect to Google OAuth2
	http.Redirect(w, r, p.Auth.AuthCodeURL(state), http.StatusFound)

	// page := pages.Home()

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

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
