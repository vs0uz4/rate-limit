package main

import (
	"log"

	"github.com/vs0uz4/rate-limit/config"
	"github.com/vs0uz4/rate-limit/internal/webserver"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := webserver.Start(cfg); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
