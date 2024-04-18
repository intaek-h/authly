package api

import (
	"net/http"
	"strconv"

	"github.com/authly/internal/authenticator"
	"github.com/authly/internal/database"
	"github.com/authly/internal/env"
	"github.com/go-chi/chi"
)

type APIs struct {
	DB   *database.Queries
	Env  env.Env
	Auth *authenticator.Authenticator
}

func (api *APIs) HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")

	userIdInt64, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		http.Error(w, "잘못된 유저 아이디입니다.", http.StatusBadRequest)
		return
	}

	user, err := api.DB.GetUserById(r.Context(), userIdInt64)
	if err != nil {
		http.Error(w, "유저를 찾을 수 없습니다.", http.StatusNotFound)
		return
	}

	var nickname *string

	if user.Nickname.Valid {
		nickname = &user.Nickname.String
	}

	w.Write([]byte(*nickname))
}
