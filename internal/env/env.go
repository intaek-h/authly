package env

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port               string `env:"PORT"`
	DatabaseUrl        string `env:"DATABASE_URL"`
	Environment        string `env:"ENVIRONMENT"`
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}

func loadEnv() (*Env, error) {
	var env Env

	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file doesn't exist.")
	}

	env.Port = os.Getenv("PORT")
	if env.Port == "" {
		return nil, errors.New("port not set")
	}

	env.DatabaseUrl = os.Getenv("DATABASE_URL")
	if env.DatabaseUrl == "" {
		return nil, errors.New("database url not set")
	}

	env.Environment = os.Getenv("ENVIRONMENT")
	if env.Environment == "" {
		env.Environment = "development"
	}

	env.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	if env.GoogleClientID == "" {
		return nil, errors.New("google client id not set")
	}

	env.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	if env.GoogleClientSecret == "" {
		return nil, errors.New("google client secret not set")
	}

	return &env, nil
}

func MustLoad() Env {
	env, err := loadEnv()
	if err != nil {
		panic(err)
	}

	return *env
}
