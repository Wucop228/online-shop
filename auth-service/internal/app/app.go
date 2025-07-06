package app

import (
	"auth-service/internal/config"
	"fmt"
)

func Run(cfg *config.Config) error {
	fmt.Println(cfg.Config_Redis.TTL)
	return nil
}
