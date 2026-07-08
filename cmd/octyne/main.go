package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/octyne/octyne/internal/app"
	"github.com/octyne/octyne/internal/config"
)

func main() {
	godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	application := app.New(cfg)

	if err := application.Server.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
