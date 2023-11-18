package api

import (
	"api-gateway/internal/client/auth_client"
	"api-gateway/internal/client/auth_client/token_client"
	"api-gateway/internal/context"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"api-gateway/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SetupRouter(ac *context.AppContext) (r *gin.Engine) {
	log.Debug().Msg("Router setup")
	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(middleware.ZerologMiddleware(log.Logger))
	r.Use(middleware.CORSMiddleware())

	authClient := auth_client.NewClient(ac.Config.HTTPServer.Other.AuthURL)
	tokenClient := token_client.NewAuthClient(authClient)

	authHandler := handler.NewHandler(tokenClient, ac.Config.HTTPServer.Other.AuthURL)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/docs/*any", authHandler.ProxyRequest("/docs/*any", nil))
			auth.GET("/accounts/me", authHandler.ProxyRequest("/accounts/me", util.StrPtr("USER")))
			auth.POST("/register", authHandler.ProxyRequest("/register", nil))
			auth.POST("/login", authHandler.ProxyRequest("/login", nil))
			auth.POST("/tokens/refresh", authHandler.ProxyRequest("/tokens/refresh", nil))
			auth.POST("/tokens/validate", authHandler.ProxyRequest("/tokens/validate", nil))
		}
	}

	log.Debug().Msg("Router setup")
	return r
}
