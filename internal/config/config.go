package config

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Configuration struct {
	HTTPServer
	Logger
}

type HTTPServer struct {
	Port  string
	Other Services
}

type Services struct {
	AuthService
	MusicService
}

type AuthService struct {
	AuthURL string
}

type MusicService struct {
	MusicFileURL     string
	MusicMetadataURL string
	MusicPlaybackURL string
}

type Logger struct {
	Level zerolog.Level
}

func LoadConfiguration() (config *Configuration, err error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config = &Configuration{
		HTTPServer{
			Port: viper.GetString("HTTP_SERVER_PORT"),
			Other: Services{
				AuthService{
					AuthURL: viper.GetString("AUTH_URL"),
				},
				MusicService{
					MusicFileURL:     viper.GetString("MUSIC_FILES_URL"),
					MusicMetadataURL: viper.GetString("MUSIC_METADATA_URL"),
					MusicPlaybackURL: viper.GetString("MUSIC_PLAYBACK_URL"),
				},
			},
		},
		Logger{
			Level: loadLoggingLevel(),
		},
	}

	return config, nil
}

func loadLoggingLevel() zerolog.Level {
	levelStr := viper.GetString("LOGGING_LEVEL")
	switch levelStr {
	case "TRACE":
		return zerolog.TraceLevel
	case "DEBUG":
		return zerolog.DebugLevel
	case "INFO":
		return zerolog.InfoLevel
	case "WARN":
		return zerolog.WarnLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "FATAL":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
