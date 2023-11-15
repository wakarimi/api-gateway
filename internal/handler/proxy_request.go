package handler

import (
	"api-gateway/internal/handler/response"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h Handler) BuildProxyUrl(endpoint string) string {
	return h.ServicePath + "/api" + endpoint
}

func (h Handler) ProxyRequest(endpoint string, role *string) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Debug().Str("servicePath", h.ServicePath).Msg("Request execution")
		url := h.BuildProxyUrl(endpoint)

		req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to create request",
				Reason:  err.Error(),
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
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to make request to the target service",
				Reason:  err.Error(),
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
