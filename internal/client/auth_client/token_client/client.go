package token_client

import "api-gateway/internal/client/auth_client"

type Client struct {
	authClient *auth_client.Client
}

func NewAuthClient(authClient *auth_client.Client) (tokenClient Client) {
	tokenClient = Client{
		authClient: authClient,
	}
	return tokenClient
}
