package config

import "github.com/caarlos0/env/v11"

type Config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`                       // постгресс подключение
	HTTPPort    string `env:"HTTP_PORT" envDefault:"8080"`                 // порт на котором  наше прилошеник пнимает запросы
	BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080"` // адрес из которого формируются короткие ссылки
	RedisAddr   string `env:"REDIS_ADDR" envDefault:"localhost:6379"`      // подключение к редис наше кеширование ссылок
}

func Load() (Config, error) {

	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
