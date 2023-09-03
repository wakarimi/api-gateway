package handlers

import (
	"api-gateway/internal/handlers/types"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func ProxyRequest(url string) func(c *gin.Context) {
	return func(c *gin.Context) {

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
