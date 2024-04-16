package page

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/authly/internal/database"
	"github.com/authly/internal/env"
	"google.golang.org/api/idtoken"
)

type Pages struct {
	DB  *database.Queries
	Env env.Env
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

func (p *Pages) HandlerOAuthPage(w http.ResponseWriter, r *http.Request) {
	env := env.MustLoad()
	code := r.URL.Query().Get("code")

	formData := url.Values{
		"client_id":     {env.GoogleClientID},
		"client_secret": {env.GoogleClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"http://localhost:8080/auth/google/callback"},
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", formData)
	if err != nil {
		http.Error(w, "요청 생성중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "구글 로그인에 실패했습니다.", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "응답 읽기중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		http.Error(w, "응답 디코딩중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	payload, err := idtoken.Validate(r.Context(), tokenResponse.IdToken, env.GoogleClientID)
	if err != nil {
		http.Error(w, "토큰 검증중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	user, err := p.DB.GetUserByEmail(r.Context(), payload.Claims["email"].(string))

	// 최초 접근이면 가입시키고 세션 생성
	if err == sql.ErrNoRows {
		user, err := p.DB.CreateUser(r.Context(), database.CreateUserParams{
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			RealName:     sql.NullString{},
			Nickname:     sql.NullString{},
			Email:        payload.Claims["email"].(string),
			ProfileImage: sql.NullString{},
		})
		if err != nil {
			http.Error(w, "유저 생성중 오류가 발생했습니다.", http.StatusInternalServerError)
			return
		}

		// 세션 생성하고, 리다이렉트
		sessionId, err := p.DB.CreateSession(r.Context(), database.CreateSessionParams{
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			http.Error(w, "세션 생성중 오류가 발생했습니다.", http.StatusInternalServerError)
			return
		}

		cookieVal := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%d", sessionId, user.ID)))

		cookie := http.Cookie{
			Name:     "session",
			Value:    cookieVal,
			Expires:  time.Now().Add(time.Hour * 24 * 7),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Secure:   true,
		}

		http.SetCookie(w, &cookie)

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if err != nil {
		http.Error(w, "유저 검색중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	sessionId, err := p.DB.CreateSession(r.Context(), database.CreateSessionParams{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		http.Error(w, "세션 생성중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	cookieVal := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%d", sessionId, user.ID)))

	cookie := http.Cookie{
		Name:     "session",
		Value:    cookieVal,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}

	// TODO: 브라우저에 쿠키가 안 감.
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
