package api

import (
	"api-gateway/internal/config"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Configuration) *gin.Engine {
	r := gin.Default()

	return r
}
