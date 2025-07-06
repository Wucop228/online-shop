package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
