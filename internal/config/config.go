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
	AuthServiceUrl string
}

func LoadConfiguration() (config *Configuration, err error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config = &Configuration{
		HttpServerConfiguration{
			Port: viper.GetString("HTTP_SERVER_PORT"),
		},
		ProxyConfiguration{
			AuthServiceUrl: viper.GetString("AUTH_SERVICE_URL"),
		},
	}

	return config, nil
}
