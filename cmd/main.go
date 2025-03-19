package main

import (
	"log"

	"github.com/imotkin/shortener/internal/config"
	"github.com/imotkin/shortener/internal/server"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	srv := server.New(cfg.Database)

	err = srv.Start(cfg.Address())
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
