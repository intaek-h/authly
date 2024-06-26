package main

import (
	"context"
	"encoding/gob"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/authly/internal/authenticator"
	"github.com/authly/internal/db"
	"github.com/authly/internal/env"
	"github.com/authly/internal/handlers/api"
	"github.com/authly/internal/handlers/page"
	m "github.com/authly/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	envr := env.MustLoad()

	db := db.MustConnect(&envr)

	auth, err := authenticator.New()
	if err != nil {
		logger.Error("Authenticator를 생성하는 중 오류가 발생했습니다.", slog.Any("error", err))
		os.Exit(1)
	}

	// 세션을 DB 대신 쿠키에 저장하는 것.
	// TODO: authentication key 를 환경변수로 관리해야함.
	store := sessions.NewCookieStore([]byte("store"))

	// TODO: 얘가 뭐하는앤지 모르겠음. 근데 이거 안하면 session.Save() 가 안됨.
	gob.Register(map[string]interface{}{})

	apis := api.APIs{
		DB:      db,
		Env:     envr,
		Auth:    auth,
		Session: store,
	}

	pages := page.Pages{
		DB:      db,
		Env:     envr,
		Auth:    auth,
		Session: store,
	}

	middlewareStore := m.MiddlewareStore{
		DB:      db,
		Env:     &envr,
		Session: store,
	}

	r := chi.NewRouter()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Pages
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.RealIP,
			middleware.RequestID,
			middleware.Logger,
			m.ContentTypeHTMLMiddleware,
			middlewareStore.HxContextMiddleware,
			middlewareStore.IsPRDContextMiddleware,
			middlewareStore.NonceMiddleware,
		)

		r.NotFound(pages.HandlerNotFoundPage)
		r.Get("/", pages.HandlerHomePage)
		r.Get("/login", pages.HandlerLoginPage)
		r.Get("/logout", pages.HandlerLogoutPage)
		r.Get("/callback", pages.HandlerOAuthPage)
	})

	// APIs
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.RealIP,
			middleware.RequestID,
			middleware.Logger,
			m.ContentTypeHTMLMiddleware,
		)

		r.Get("/users/{id}", apis.HandlerGetUser)
	})

	killSignal := make(chan os.Signal, 1)

	signal.Notify(killSignal, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:        ":" + envr.Port,
		Handler:     r,
		ReadTimeout: time.Second * 10,
	}

	go func() {
		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("서버가 정상적으로 종료되었습니다.")
		} else if err != nil {
			logger.Error("서버가 비정상적으로 종료되었습니다.", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	logger.Info("서버가 시작되었습니다.", slog.String("포트", envr.Port))

	// 종료 시그널을 기다립니다.
	<-killSignal

	// 종료 시그널을 받으면 5초 뒤에 서버를 종료합니다.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("서버 종료 중 오류가 발생했습니다.", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("서버가 종료되었습니다.")
}
