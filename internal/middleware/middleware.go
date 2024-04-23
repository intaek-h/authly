package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/authly/internal/database"
	"github.com/authly/internal/env"
	"github.com/gorilla/sessions"
)

type MiddlewareStore struct {
	DB      *database.Queries
	Env     *env.Env
	Session *sessions.CookieStore
}

func ContentTypeHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

type UserContextKey string

var UserKey UserContextKey = "user"

// TODO: 필요 없음. 지우기.
func (md *MiddlewareStore) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		decodedValue, err := base64.StdEncoding.DecodeString(sessionCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		splitValue := strings.Split(string(decodedValue), ":")

		if len(splitValue) != 2 {
			next.ServeHTTP(w, r)
			return
		}

		sessionID := splitValue[0]

		sid, err := strconv.ParseInt(sessionID, 10, 64)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := md.DB.GetUserFromSession(r.Context(), sid)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if time.Now().After(user.SessionExpiresAt) {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, &user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *database.GetUserFromSessionRow {
	user := ctx.Value(UserKey)
	if user == nil {
		return nil
	}

	return user.(*database.GetUserFromSessionRow)
}

func (md *MiddlewareStore) StateCookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// session, _ := md.Session.Get(r, "auth-session")

		// if !session.IsNew {
		// 	profile := session.Values["profile"].(map[string]interface{})
		// 	fmt.Printf("profile: %#v\n", profile["nickname"])
		// }

		next.ServeHTTP(w, r)
	})
}

func (md *MiddlewareStore) HxContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHxRequest := r.Header.Get("HX-Request") == "true"

		ctx := context.WithValue(r.Context(), "isHxRequest", isHxRequest)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (md *MiddlewareStore) IsPRDContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isProduction := md.Env.Environment == "production"

		ctx := context.WithValue(r.Context(), "isProduction", isProduction)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

type key string

var NonceKey key = "nonce"
var nonce = generateRandomString(16)

// Refused to apply inline style because it violates the following Content Security Policy directive: "default-src 'self'".
// Either the 'unsafe-inline' keyword, a hash ('sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg='),
// or a nonce ('nonce-...') is required to enable inline execution.
// Note also that 'style-src' was not explicitly set, so 'default-src' is used as a fallback.
var htmxInlineStyleCSPHash = "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg="

// TODO: CSP 에 대해 공부할 필요 있음. nonce 하나로 공유하는게 맞나??
func (md *MiddlewareStore) NonceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), NonceKey, nonce)

		// insert the nonces into the content security policy header
		cspHeader := fmt.Sprintf("default-src 'self'; script-src 'nonce-%s' ; style-src 'nonce-%s' '%s' ;", nonce, nonce, htmxInlineStyleCSPHash)

		w.Header().Set("Content-Security-Policy", cspHeader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
