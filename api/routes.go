package api

import (
	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Configuration) *gin.Engine {
	authBasePath := cfg.AuthServiceUrl + "/api/auth-service"

	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.ProxyRequest(authBasePath+"/register"))
			auth.POST("/login", handlers.ProxyRequest(authBasePath+"/login"))
			auth.POST("/refresh", handlers.ProxyRequest(authBasePath+"/refresh"))
			auth.POST("/validate", handlers.ProxyRequest(authBasePath+"/validate"))
		}
	}

	return r
}
