package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

type pageResources struct{}

type Data struct {
	GoogleAuthClientId    string
	GoogleAuthScope       string
	GoogleAuthRedirectUri string
}

func (pg pageResources) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", pg.Home)
	r.Get("/auth/google/callback", pg.GoogleCallback)

	return r
}

func (pg pageResources) Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("browser/pages/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := Data{
		GoogleAuthClientId:    os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleAuthScope:       "https://www.googleapis.com/auth/cloud-platform",
		GoogleAuthRedirectUri: "http://localhost:8080/api/v1/users/oauth2callback",
	}

	tmpl := template.Must(t, err)
	tmpl.Execute(w, data)
}

func (pg pageResources) GoogleCallback(w http.ResponseWriter, r *http.Request) {

	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			log.Println(name, value)
		}
	}

	scope := r.URL.Query().Get("code")
	log.Println(scope)

	t, err := template.ParseFiles("browser/pages/login_success.html")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// http://localhost:8080/auth/google/callback
	// #access_token=ya29.a0Ad52N3-tEiP5EQ2vjpidmTKX9dNmcaQZxyt66AQ0o9qhajblzLRT4ww8pREQc0M74EgWno9pL7vemHWSLXh1PhVtiwJSx0XdHtY0Z8sQgerZY5Iv1xawOFIl-lJ0JevH1aehVUV9uQugSrPsfymFfBiMIWVttXZD-AaCgYKAcASARMSFQHGX2Mi92N3IQsQahVMsZgmgs4ffw0169
	// &token_type=Bearer
	// &expires_in=3599
	// &scope=https://www.googleapis.com/auth/cloud-platform
	tmpl := template.Must(t, err)
	tmpl.Execute(w, nil)
}
