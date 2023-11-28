package context

import (
	"api-gateway/internal/config"
)

type AppContext struct {
	Config *config.Configuration
}
