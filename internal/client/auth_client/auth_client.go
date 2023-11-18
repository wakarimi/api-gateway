package auth_client

import (
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Client struct {
	BaseUrl    string
	HttpClient *http.Client
}

func NewClient(baseUrl string) (client *Client) {
	client = &Client{
		BaseUrl: baseUrl,
		HttpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	return client
}

func (c *Client) Request(method, path string, body io.Reader) (*http.Response, error) {
	reqURL := c.BaseUrl + path

	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		log.Error().Err(err).Str("method", method).Str("path", path).Msg("Failed to create request")
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("method", method).Str("path", path).Msg("Failed to execute request")
		return nil, err
	}

	return resp, nil
}
