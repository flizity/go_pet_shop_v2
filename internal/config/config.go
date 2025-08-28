package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	DatabaseURL string `yaml:"database_url" env:"DATABASE_URL"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:5434"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

const configPath = "./config/local.yaml"

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("error loading .env file")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config " + err.Error())
	}

	if cfg.DatabaseURL = os.Getenv("DATABASE_URL"); cfg.DatabaseURL == "" {
		panic("error loading DATABASE_URL from environment, please set it in .env file")
	}

	return &cfg
}
