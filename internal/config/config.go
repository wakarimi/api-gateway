package config

import (
	"github.com/spf13/viper"
	"strings"
)

type Configuration struct {
	HttpServerConfiguration
	ProxyConfiguration
}

type HttpServerConfiguration struct {
	Port string
}

type ProxyConfiguration struct {
	Services map[string]ServiceConfig
}

type ServiceConfig struct {
	BaseUrl    string
	PathPrefix string
}

func LoadConfiguration() (config *Configuration, err error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config = &Configuration{
		HttpServerConfiguration{
			Port: viper.GetString("HTTP_SERVER_PORT"),
		},
		ProxyConfiguration{
			Services: map[string]ServiceConfig{
				"auth": {
					BaseUrl:    viper.GetString("AUTH_SERVICE_URL"),
					PathPrefix: "/api/auth-service",
				},
			},
		},
	}

	return config, nil
}
