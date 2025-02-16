package main

import (
	"context"
	"github.com/nglmq/avito-shop/internal/app/transaction"
	"github.com/nglmq/avito-shop/internal/config"
	"github.com/nglmq/avito-shop/internal/storage/postgresql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nglmq/avito-shop/internal/api/handlers"
	"github.com/nglmq/avito-shop/internal/app/auth"
	"github.com/nglmq/avito-shop/internal/app/history"
	"github.com/nglmq/avito-shop/internal/app/merch"
	md "github.com/nglmq/avito-shop/internal/middleware"
)

func main() {
	config.ParseFlags()

	storage, _ := postgresql.NewRepo(context.Background(), config.DatabaseDSN)
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	authService := auth.New(logger, storage)
	infoService := history.New(logger, storage)
	txService := transaction.New(logger, storage)
	merchService := merch.New(logger, storage)

	router := chi.NewRouter()
	router.Use(middleware.DefaultLogger)

	authMiddleware := md.CheckAuthMiddleware(logger)
	router.Route("/api/", func(r chi.Router) {
		r.With(authMiddleware).Get("/info", handlers.HandleGetInfo(infoService))
		r.With(authMiddleware).Get("/buy/{item}", handlers.HandleBuyItem(merchService))
		r.Post("/auth", handlers.HandleAuth(authService))
		r.With(authMiddleware).Post("/sendCoin", handlers.HandleSendCoin(txService))
	})

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 2 * time.Second,
	}

	logger.Info("Starting server on port :8080")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
