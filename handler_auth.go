package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type authResources struct{}

func (ar authResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/google/callback", ar.GoogleCallback)

	return r
}

func (ar authResources) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// get all query params
	query := r.URL.Query()
	for k, v := range query {
		log.Printf("%s: %s", k, v)
	}

	w.Write([]byte("google callback"))
}
