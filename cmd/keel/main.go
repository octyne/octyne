package main

import (
	"log"

	"github.com/usekeel/keel/internal/app"
	"github.com/usekeel/keel/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	application := app.New(cfg)

	if err := application.Server.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
