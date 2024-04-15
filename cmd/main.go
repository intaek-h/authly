package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/authly/internal/env"
	"github.com/authly/internal/handlers/page"
	m "github.com/authly/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	envr := env.MustLoad()

	// db := db.MustConnect(envr.DatabaseUrl)

	r := chi.NewRouter()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.RealIP,
			middleware.RequestID,
			middleware.Logger,
			// middleware.Recoverer,
			m.TextHTMLMiddleware,
		)

		r.NotFound(page.NewNotFoundPageHandler().ServeHTTP)

		r.Get("/", page.NewHomePageHandler().ServeHTTP)
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
