package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

type usersResource struct{}

func (ur usersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", ur.Create)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", ur.Get)
	})

	return r
}

func (ur usersResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get user"))
}

func (ur usersResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create user"))
}
