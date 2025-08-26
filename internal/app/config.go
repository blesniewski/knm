package app

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr string `env:"LISTEN_ADDR" envDefault:":8080"`
	Orx        struct {
		BaseURL string `env:"BASE_URL"        envDefault:"https://openexchangerates.org/api"`
		AppID   string `env:"APP_ID,required"`
	} `envPrefix:"OPENEXCHANGERATES_"`
}

func NewConfig() (Config, error) {
	_ = godotenv.Load()

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
