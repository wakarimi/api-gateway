package main

import (
	"api-gateway/api"
	"api-gateway/internal/config"
	"log"
)

// @title Wakarimi ApiGateway API
// @version 0.1
// @description This is the api gateway service for Wakarimi.
// @contact.name Zalimannard
// @contact.email zalimannard@mail.ru
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8021
// @BasePath /api/api-gateway-service
func main() {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		log.Fatal("Failed to load configuration: %v", err)
	}

	r := api.SetupRouter(cfg)
	r.Run(":" + cfg.Port)
}
