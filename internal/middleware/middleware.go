package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/authly/internal/database"
)

type MiddlewareStore struct {
	DB *database.Queries
}

func ContentTypeHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// TODO: HX-Request 여부 처리하는 미들웨어 구현.
// func HxRequestMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Header.Get("HX-Request") == "true" {
// 			r.Context().Value("hx-request")
// 		}
// 	})
// }

type UserContextKey string

var UserKey UserContextKey = "user"

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
		stateCookie, _ := r.Cookie("state")

		fmt.Printf("stateCookie: %v\n", stateCookie)

		next.ServeHTTP(w, r)
	})
}
