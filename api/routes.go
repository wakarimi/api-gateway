package api

import (
	"api-gateway/internal/context"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))
	r.Use(middleware.CORSMiddleware())

	authHandler := handler.NewHandler(ac.Config.HTTPServer.Other.AuthURL)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/docs/*any", authHandler.ProxyRequest("/docs/*any", nil))
			auth.POST("/register", authHandler.ProxyRequest("/register", nil))
			auth.POST("/login", authHandler.ProxyRequest("/login", nil))
			auth.POST("/tokens/refresh", authHandler.ProxyRequest("/tokens/refresh", nil))
			auth.POST("/tokens/validate", authHandler.ProxyRequest("/tokens/validate", nil))
		}
	}

	log.Debug().Msg("Router setup")
	return r
}
