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
	metadataHandler := handler.NewHandler(tokenClient, ac.Config.HTTPServer.Other.MusicMetadataURL)

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

		metadata := api.Group("/metadata")
		{
			metadata.GET("/docs/*any", metadataHandler.ProxyRequest("/docs/*any", nil))

			metadata.GET("/albums", metadataHandler.ProxyRequest("/albums", util.StrPtr("USER")))
			metadata.GET("/albums/:albumId", metadataHandler.ProxyRequest("/albums/:albumId", util.StrPtr("USER")))
			metadata.GET("/albums/:albumId/songs", metadataHandler.ProxyRequest("/albums/:albumId/songs", util.StrPtr("USER")))
			metadata.GET("/albums/:albumId/covers", metadataHandler.ProxyRequest("/albums/:albumId/covers", util.StrPtr("USER")))

			metadata.GET("/artists", metadataHandler.ProxyRequest("/artists", util.StrPtr("USER")))
			metadata.GET("/artists/:artistId", metadataHandler.ProxyRequest("/artists/:artistId", util.StrPtr("USER")))
			metadata.GET("/artists/:artistId/songs", metadataHandler.ProxyRequest("/artists/:artistId/songs", util.StrPtr("USER")))
			metadata.GET("/artists/:artistId/covers", metadataHandler.ProxyRequest("/artists/:artistId/covers", util.StrPtr("USER")))

			metadata.GET("/genres", metadataHandler.ProxyRequest("/genres", util.StrPtr("USER")))
			metadata.GET("/genres/:genreId", metadataHandler.ProxyRequest("/genres/:genreId", util.StrPtr("USER")))
			metadata.GET("/genres/:genreId/songs", metadataHandler.ProxyRequest("/genres/:genreId/songs", util.StrPtr("USER")))
			metadata.GET("/genres/:genreId/covers", metadataHandler.ProxyRequest("/genres/:genreId/covers", util.StrPtr("USER")))
		}
	}

	log.Debug().Msg("Router setup")
	return r
}
