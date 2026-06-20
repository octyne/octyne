package main

import (
	"log"

	"github.com/usekeel/keel/internal/config"
	"github.com/usekeel/keel/internal/server"
)

func main() {
	cfg := config.Load()
	srv := server.New()

	if err := srv.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
