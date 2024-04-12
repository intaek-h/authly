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

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", pageResources{}.Routes())
	r.Mount("/js", jsResources{}.Routes())
	r.Mount("/api/v1/", v1)
	v1.Mount("/users", usersResource{}.Routes())
	v1.Mount("/auth", authResources{}.Routes())

	v1.Get("/healthz", handlerReady)

	log.Printf("server is on port %s", port)
	http.ListenAndServe(":"+port, r)
}
