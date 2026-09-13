package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abhisheksinha-989/ReCircle/internal/config"
	"github.com/abhisheksinha-989/ReCircle/internal/db"
	"github.com/abhisheksinha-989/ReCircle/internal/handlers"
)

func main() {

	cfg := config.MustLoad()

	db, err := db.Connectdb(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Database: %v", err)
	}

	fmt.Println("database opened")
	fmt.Println("starting serverr")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /check_health", handlers.CheckHealth)
	mux.HandleFunc("GET /listings", handlers.List(db))

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
