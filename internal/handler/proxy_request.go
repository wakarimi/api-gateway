package handler

import (
	"api-gateway/internal/handler/response"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h Handler) BuildProxyUrl(endpoint string) string {
	return h.ServicePath + "/api" + endpoint
}

func (h Handler) ProxyRequest(endpoint string, role *string) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Debug().Str("servicePath", h.ServicePath).Msg("Request execution")
		params := c.Params
		queryParams := c.Request.URL.Query()

		url := h.BuildProxyUrl(endpoint)

		for _, p := range params {
			url = strings.Replace(url, ":"+p.Key, p.Value, 1)
		}
		if len(queryParams) > 0 {
			url += "?" + queryParams.Encode()
		}

		req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Error{
				Message: "Failed to create request",
				Reason:  err.Error(),
			})
			return
		}

		log.Debug().Msg("Fetching account details")
		if role != nil {
			authTokenWithType := c.GetHeader("Authorization")
			if authTokenWithType == "" {
				err := fmt.Errorf("failed to get authorization header")
				c.JSON(http.StatusUnauthorized, response.Error{
					Message: "Failed to get Authorization header",
					Reason:  err.Error(),
				})
				return
			}

			if !strings.HasPrefix(authTokenWithType, "Bearer ") {
				err := fmt.Errorf("invalid auth token format")
				c.JSON(http.StatusUnauthorized, response.Error{
					Message: "Invalid auth token format",
					Reason:  err.Error(),
				})
				return
			}

			authToken := strings.TrimPrefix(authTokenWithType, "Bearer ")

			validateResponse, err := h.TokenClient.Validate(authToken)
			if err != nil {
				log.Error().Err(err).Msg("Failed to validate token")
				c.JSON(http.StatusUnauthorized, response.Error{
					Message: "Failed to check auth token",
					Reason:  err.Error(),
				})
				return
			}

			if validateResponse.Valid {
				accountID := strconv.Itoa(*validateResponse.AccountID)
				deviceID := strconv.Itoa(*validateResponse.DeviceID)
				req.Header.Add("X-Account-ID", accountID)
				req.Header.Add("X-Device-ID", deviceID)
			} else {
				err := fmt.Errorf("invalid token")
				log.Error().Err(err).Msg("Invalid token")
				c.JSON(http.StatusUnauthorized, response.Error{
					Message: "Invalid token",
					Reason:  err.Error(),
				})
				return
			}
		}

		for name, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to make request to the target service")
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
