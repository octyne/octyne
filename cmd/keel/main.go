package main

import (
	"log"

	"github.com/usekeel/keel/internal/app"
	"github.com/usekeel/keel/internal/config"
)

func main() {
	cfg := config.Load()
	application := app.New()

	if err := application.Server.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
