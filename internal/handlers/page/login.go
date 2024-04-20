package page

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func (p *Pages) HandlerLoginPage(w http.ResponseWriter, r *http.Request) {
	sessionId, err := generateRandomState()
	if err != nil {
		http.Error(w, "로그인 상태를 생성하는 중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	session, _ := p.Session.Get(r, "auth-session")
	session.Values["state"] = sessionId

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "세션 저장 중 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	// Redirect to Google OAuth2
	http.Redirect(w, r, p.Auth.AuthCodeURL(sessionId), http.StatusFound)
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
