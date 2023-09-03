package main

import (
	"api-gateway/api"
	"api-gateway/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		log.Fatal("Failed to load configuration: %v", err)
	}

	r := api.SetupRouter(cfg)
	r.Run(":" + cfg.Port)
}
