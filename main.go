package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file doesn't exist.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port not set.")
	}

	r := chi.NewRouter()
	v1 := chi.NewRouter()

	r.Mount("/v1", v1)

	v1.Use(middleware.RealIP)
	v1.Use(middleware.RequestID)
	v1.Use(middleware.Logger)
	v1.Use(middleware.Recoverer)

	v1.Get("/healthz", handlerReady)

	v1.Mount("/users", usersResource{}.Routes())

	log.Printf("server is on port %s", port)
	http.ListenAndServe(":"+port, r)
}
