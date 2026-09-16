package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/abhisheksinha-989/ReCircle/internal/config"
	"github.com/abhisheksinha-989/ReCircle/internal/db"
	"github.com/abhisheksinha-989/ReCircle/internal/handlers"
	"github.com/abhisheksinha-989/ReCircle/internal/middleware"
)

func main() {

	cfg := config.MustLoad()

	db, err := db.Connectdb(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Database: %v", err)
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	fmt.Println("database opened")
	fmt.Println("starting serverr")

	lh := handlers.NewListingHandler(db, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /check_health", handlers.CheckHealth)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	handler := middleware.RequestId(mux)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
