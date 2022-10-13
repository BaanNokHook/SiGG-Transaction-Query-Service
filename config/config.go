package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App
		HTTP
		Log
		REDIS
	}

	// App -.
	App struct {
		Name    string `env-required:"true" env:"APP_NAME"`
		Version string `env-required:"true" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	REDIS struct {
		Addr     string `env-required:"true" env:"REDIS_ADDR"`
		Password string `env-required:"true" env:"REDIS_PASSWORD"`
		DB       int    `env-required:"true" env:"REDIS_DB"`
	}

	// Pusher struct {
	// 	AppID   string `env-required:"false" env:"APPID"`
	// 	Key     string `env-required:"false" env:"KEY"`
	// 	Secret  string `env-required:"false" env:"SECRET"`
	// 	Cluster string `env-required:"false" env:"CLUSTER"`
	// 	Secure  bool   `env-required:"false" env:"SECURE"`
	// }
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(".env"); err == nil {
		err = cleanenv.ReadConfig(".env", cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
