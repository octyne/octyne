package main

import (
	"github.com/usekeel/keel/internal/server"
	"log"
)

func main() {
	srv := server.New()

	if err := srv.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
