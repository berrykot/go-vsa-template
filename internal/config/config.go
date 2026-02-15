package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Env  string `env:"APP_ENV" envDefault:"development"`
	Port int    `env:"PORT" envDefault:"8080"`
	//Database struct {
	//	Host string `env:"DB_HOST,required"`
	//	Port int `env:"DB_PORT,required"`
	//}
	Cron struct {
		HealthCron string `env:"HEALTH_CRON,required"`
		//OtherСron string `env:"OTHER_СRON,required"`
	}
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("Load env: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
