package handler

import "api-gateway/internal/client/auth_client/token_client"

type Handler struct {
	TokenClient token_client.Client
	ServicePath string
}

func NewHandler(tokenClient token_client.Client,
	servicePath string,
) (h *Handler) {
	h = &Handler{
		TokenClient: tokenClient,
		ServicePath: servicePath,
	}

	return h
}
