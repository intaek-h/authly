package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

type jsResources struct{}

func (js jsResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/{filename}", js.Get)

	return r
}

func (js jsResources) Get(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	http.ServeFile(w, r, "browser/js/"+filename)
}
