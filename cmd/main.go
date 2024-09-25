package main

import (
	"github.com/imotkin/shortener/pkg/config"
	"github.com/imotkin/shortener/pkg/server"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	srv := server.New(cfg.Database)

	err = srv.Start(cfg.Address())
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
