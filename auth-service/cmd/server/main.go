package main

import (
	"github.com/Wucop228/online-shop/auth-service/internal/app"
	"github.com/Wucop228/online-shop/auth-service/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
