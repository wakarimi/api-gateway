package main

import (
	"api-gateway/api"
	"api-gateway/internal/config"
	"api-gateway/internal/context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title Wakarimi API-Gateway
// @version 0.5

// @contact.name Dmitry Kolesnikov (Zalimannard)
// @contact.email zalimannard@mail.ru

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8021
// @BasePath /api
func main() {
	cfg := loadConfiguration()

	initializeLogger(cfg.Logger.Level)

	ctx := context.AppContext{
		Config: cfg,
	}

	server := initializeServer(&ctx)
	runServer(server, cfg.HTTPServer.Port)
}

func loadConfiguration() *config.Configuration {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to load configuration")
	}
	log.Info().Msg("Configuration loaded")
	return cfg
}

func initializeLogger(level zerolog.Level) (logger *zerolog.Logger) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().Caller().Logger().
		With().Str("service", "authentication").Logger().
		Level(level)
	log.Debug().Msg("Logger initialized")
	return &log.Logger
}

func initializeServer(ctx *context.AppContext) (r *gin.Engine) {
	r = api.SetupRouter(ctx)
	log.Debug().Msg("Router initialized")
	return r
}

func runServer(server *gin.Engine, port string) {
	if err := server.Run(":" + port); err != nil {
		log.Panic().Err(err).Msg("Failed to start server")
	}
}
