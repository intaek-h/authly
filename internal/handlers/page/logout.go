package page

import (
	"net/http"
	"net/url"
	"os"
)

func (p *Pages) HandlerLogoutPage(w http.ResponseWriter, r *http.Request) {
	session, err := p.Session.Get(r, "auth-session")
	if session.IsNew || err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Values["profile"] = nil
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "세션 초기화 중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		http.Error(w, "로그아웃 URL을 생성하는 중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host)
	if err != nil {
		http.Error(w, "로그아웃 URL을 생성하는 중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))

	logoutUrl.RawQuery = parameters.Encode()

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
