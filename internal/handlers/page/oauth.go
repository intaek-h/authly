package page

import (
	"fmt"
	"net/http"

	"github.com/authly/internal/authenticator"
	"github.com/authly/internal/database"
	"github.com/authly/internal/env"
	"github.com/gorilla/sessions"
)

type Pages struct {
	DB      *database.Queries
	Env     env.Env
	Auth    *authenticator.Authenticator
	Session *sessions.CookieStore
}

type TokenPayload struct {
	Issuer   string `json:"Issuer"`
	Audience string `json:"Audience"`
	Expires  int    `json:"Expires"`
	IssuedAt int    `json:"IssuedAt"`
	Subject  string `json:"Subject"`
	Claims   struct {
		AtHash        string  `json:"at_hash"`
		Aud           string  `json:"aud"`
		Azp           string  `json:"azp"`
		Email         string  `json:"email"`
		EmailVerified bool    `json:"email_verified"`
		Exp           float64 `json:"exp"`
		FamilyName    string  `json:"family_name"`
		GivenName     string  `json:"given_name"`
		Iat           float64 `json:"iat"`
		Iss           string  `json:"iss"`
		Name          string  `json:"name"`
		Picture       string  `json:"picture"`
		Sub           string  `json:"sub"`
	} `json:"Claims"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}

//	{
//		"aud":"FQ2KsfdsfdbLRvOFhyqcGoj",
//		"exp":1.713v31e+09,
//		"family_name":"Hwang",
//		"given_name":"Intaek",
//		"iat":1.714431e+09,
//		"iss":"https://dev-l0hgsdfdercu0due3d.jp.auth0.com/",
//		"locale":"ko",
//		"name":"Intaek Hwang",
//		"nickname":"a4fou2ahiou",
//		"picture":"https://lh3.googleusercontent.com/a/ACg8ocJcT_z9RY423EqN52DBa9xsSobcKpyaN4RXg5-1L11ujNBg=s96-c",
//		"sid":"NAfds14Ucq7aka8MDAtKm5Mw2qi",
//		"sub":"google-oauth2|112234219858648698",
//	}
type GoogleProfile struct {
	Aud        string `json:"aud"`
	Exp        int    `json:"exp"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Iat        int    `json:"iat"`
	Iss        string `json:"iss"`
	Locale     string `json:"locale"`
	Name       string `json:"name"`
	Nickname   string `json:"nickname"`
	Picture    string `json:"picture"`
	Sid        string `json:"sid"`
	Sub        string `json:"sub"`
}

func (p *Pages) HandlerOAuthPage(w http.ResponseWriter, r *http.Request) {
	session, err := p.Session.Get(r, "auth-session")
	if session.IsNew || err != nil {
		http.Error(w, "세션스토어에 접근할 수 없습니다", http.StatusInternalServerError)
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "State is missing.", http.StatusBadRequest)
		return
	}

	if session.Values["state"] != r.URL.Query().Get("state") {
		http.Error(w, "Invalid state parameter.", http.StatusBadRequest)
		return
	}

	// Exchange an authorization code for a token.
	code := r.URL.Query().Get("code")
	token, err := p.Auth.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token.", http.StatusInternalServerError)
		return
	}

	idToken, err := p.Auth.VerifyIDToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to verify ID token.", http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, "Failed to parse ID token claims.", http.StatusInternalServerError)
		return
	}
	fmt.Printf("profile: %#v\n", profile)
	session.Values["profile"] = profile
	session.Values["access_token"] = token.AccessToken

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session.", http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page.
	referer := r.Referer()
	if referer == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, referer, http.StatusTemporaryRedirect)
}
