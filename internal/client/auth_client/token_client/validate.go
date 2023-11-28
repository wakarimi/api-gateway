package token_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type validateRequest struct {
	AccessToken string `json:"accessToken"`
}

type validateResponse struct {
	Valid     bool      `json:"valid"`
	AccountID *int      `json:"accountId"`
	DeviceID  *int      `json:"deviceId"`
	Roles     *[]string `json:"roles"`
	IssuedAt  *int64    `json:"issuedAt"`
	ExpiryAt  *int64    `json:"expiryAt"`
}

func (c *Client) Validate(accessToken string) (validationResponse validateResponse, err error) {
	log.Debug().Msg("Fetching validation info")

	requestBody := validateRequest{
		AccessToken: accessToken,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal requestBody")
		return validateResponse{}, err
	}

	resp, err := c.authClient.Request(http.MethodPost, "/api/tokens/validate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute request")
		return validateResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("received unexpected status code: %d", resp.StatusCode)
		log.Error().Err(err).Str("statusCode", resp.Status).Msg("Received unexpected status code")
		return validateResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return validateResponse{}, err
	}

	err = json.Unmarshal(body, &validationResponse)
	if err != nil {
		log.Error().Err(err).Msg("Failed to deserialize response body")
		return validateResponse{}, err
	}

	return validationResponse, nil
}
