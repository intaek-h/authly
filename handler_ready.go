package main

import "net/http"

func handlerReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
