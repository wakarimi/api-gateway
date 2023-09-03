package handlers

import (
	"api-gateway/internal/config"
	"api-gateway/internal/handlers/types"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func BuildProxyUrl(cfg *config.Configuration, serviceKey string, originalPath string) string {
	if serviceConfig, exists := cfg.ProxyConfiguration.Services[serviceKey]; exists {
		return serviceConfig.BaseUrl + serviceConfig.PathPrefix + strings.TrimPrefix(originalPath, "/api/"+serviceKey)
	}
	return ""
}

func ProxyRequest(cfg *config.Configuration, serviceKey string) func(c *gin.Context) {
	return func(c *gin.Context) {
		url := BuildProxyUrl(cfg, serviceKey, c.Request.URL.Path)

		req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.Error{
				Error: "Failed to create request",
			})
			return
		}

		for name, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.Error{
				Error: "Failed to make request to the target service",
			})
			return
		}
		defer resp.Body.Close()

		for name, values := range resp.Header {
			for _, value := range values {
				c.Header(name, value)
			}
		}
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}
