package main

import (
	"github.com/timurkash/queue2/internal/app"
	"log"
)

func main() {
	cfg := app.ParseFlags()
	application := app.New(cfg)

	log.Printf("Starting queue2 on port %d", cfg.Port)
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}

}
