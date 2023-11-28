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
	musicFilesHandler := handler.NewHandler(tokenClient, ac.Config.HTTPServer.Other.MusicFileURL)
	musicMetadataHandler := handler.NewHandler(tokenClient, ac.Config.HTTPServer.Other.MusicMetadataURL)

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

		musicFiles := api.Group("/music-files")
		{
			musicFiles.GET("/docs/*any", musicFilesHandler.ProxyRequest("/docs/*any", nil))

			musicFiles.GET("/roots", musicFilesHandler.ProxyRequest("/roots", util.StrPtr("USER")))
			musicFiles.POST("/roots", musicFilesHandler.ProxyRequest("/roots", util.StrPtr("ADMIN")))
			musicFiles.DELETE("/roots/:dirId", musicFilesHandler.ProxyRequest("/roots/:dirId", util.StrPtr("ADMIN")))

			musicFiles.GET("/dirs/:dirId", musicFilesHandler.ProxyRequest("/dirs/:dirId", util.StrPtr("USER")))
			musicFiles.GET("/dirs/:dirId/content", musicFilesHandler.ProxyRequest("/dirs/:dirId/content", util.StrPtr("USER")))
			musicFiles.POST("/dirs/:dirId/scan", musicFilesHandler.ProxyRequest("/dirs/:dirId/scan", util.StrPtr("ADMIN")))
			musicFiles.POST("/dirs/scan", musicFilesHandler.ProxyRequest("/dirs/scan", util.StrPtr("ADMIN")))

			musicFiles.GET("/audio-files/:audioFileId", musicFilesHandler.ProxyRequest("/audio-files/:audioFileId", util.StrPtr("USER")))
			musicFiles.GET("/audio-files", musicFilesHandler.ProxyRequest("/audio-files", util.StrPtr("USER")))
			musicFiles.GET("/audio-files/:audioFileId/download", musicFilesHandler.ProxyRequest("/audio-files/:audioFileId/download", util.StrPtr("USER")))
			musicFiles.GET("/audio-files/:audioFileId/cover", musicFilesHandler.ProxyRequest("/audio-files/:audioFileId/cover", util.StrPtr("USER")))
			musicFiles.GET("/audio-files/sha256/:sha256", musicFilesHandler.ProxyRequest("/audio-files/sha256/:sha256", util.StrPtr("USER")))
			musicFiles.PUT("/audio-files/covers-top", musicFilesHandler.ProxyRequest("/audio-files/covers-top", util.StrPtr("USER")))

			musicFiles.GET("/covers/:coverId", musicFilesHandler.ProxyRequest("/covers/:coverId", util.StrPtr("USER")))
			musicFiles.GET("/covers/:coverId/image", musicFilesHandler.ProxyRequest("/covers/:coverId/image", nil))
		}

		musicMetadata := api.Group("/music-metadata")
		{
			musicMetadata.GET("/docs/*any", musicMetadataHandler.ProxyRequest("/docs/*any", nil))

			musicMetadata.GET("/songs", musicMetadataHandler.ProxyRequest("/songs", util.StrPtr("USER")))
			musicMetadata.GET("/songs/:songId", musicMetadataHandler.ProxyRequest("/songs/:songId", util.StrPtr("USER")))

			musicMetadata.GET("/albums", musicMetadataHandler.ProxyRequest("/albums", util.StrPtr("USER")))
			musicMetadata.GET("/albums/:albumId", musicMetadataHandler.ProxyRequest("/albums/:albumId", util.StrPtr("USER")))
			musicMetadata.GET("/albums/:albumId/songs", musicMetadataHandler.ProxyRequest("/albums/:albumId/songs", util.StrPtr("USER")))
			musicMetadata.GET("/albums/:albumId/covers", musicMetadataHandler.ProxyRequest("/albums/:albumId/covers", util.StrPtr("USER")))

			musicMetadata.GET("/artists", musicMetadataHandler.ProxyRequest("/artists", util.StrPtr("USER")))
			musicMetadata.GET("/artists/:artistId", musicMetadataHandler.ProxyRequest("/artists/:artistId", util.StrPtr("USER")))
			musicMetadata.GET("/artists/:artistId/songs", musicMetadataHandler.ProxyRequest("/artists/:artistId/songs", util.StrPtr("USER")))
			musicMetadata.GET("/artists/:artistId/covers", musicMetadataHandler.ProxyRequest("/artists/:artistId/covers", util.StrPtr("USER")))

			musicMetadata.GET("/genres", musicMetadataHandler.ProxyRequest("/genres", util.StrPtr("USER")))
			musicMetadata.GET("/genres/:genreId", musicMetadataHandler.ProxyRequest("/genres/:genreId", util.StrPtr("USER")))
			musicMetadata.GET("/genres/:genreId/songs", musicMetadataHandler.ProxyRequest("/genres/:genreId/songs", util.StrPtr("USER")))
			musicMetadata.GET("/genres/:genreId/covers", musicMetadataHandler.ProxyRequest("/genres/:genreId/covers", util.StrPtr("USER")))
		}
	}

	log.Debug().Msg("Router setup")
	return r
}
