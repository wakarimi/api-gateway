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

type ValidateResponse struct {
	Valid     bool      `json:"valid"`
	AccountID *int      `json:"accountId"`
	DeviceID  *int      `json:"deviceId"`
	Roles     *[]string `json:"roles"`
	IssuedAt  *int64    `json:"issuedAt"`
	ExpiryAt  *int64    `json:"expiryAt"`
}

func (c *Client) Validate(accessToken string) (validationResponse ValidateResponse, err error) {
	log.Debug().Msg("Fetching validation info")

	requestBody := validateRequest{
		AccessToken: accessToken,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal requestBody")
		return ValidateResponse{}, err
	}

	resp, err := c.authClient.Request(http.MethodPost, "/api/tokens/verify", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute request")
		return ValidateResponse{}, err
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
		return ValidateResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return ValidateResponse{}, err
	}

	err = json.Unmarshal(body, &validationResponse)
	if err != nil {
		log.Error().Err(err).Msg("Failed to deserialize response body")
		return ValidateResponse{}, err
	}

	return validationResponse, nil
}
