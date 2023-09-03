package api

import (
	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Configuration) *gin.Engine {

	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/docs/*any", handlers.ProxyRequest(cfg, "auth"))
			auth.POST("/register", handlers.ProxyRequest(cfg, "auth"))
			auth.POST("/login", handlers.ProxyRequest(cfg, "auth"))
			auth.POST("/refresh", handlers.ProxyRequest(cfg, "auth"))
			auth.POST("/validate", handlers.ProxyRequest(cfg, "auth"))
		}
	}

	return r
}
