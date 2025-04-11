package config

import (
    "flag"
    "log"
    "os"

    "github.com/caarlos0/env/v11"
)

type Config struct {
    Addr string `env:"ADDRESS" envDefault:"localhost:8080"`
}

// Метод для получения адреса
func (c *Config) GetAddr() string {
    return c.Addr
}

func NewConfig() (*Config, error) {
    cfg := &Config{}

    // Парсим переменные окружения
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }

    // Устанавливаем флаг командной строки
    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    flagSet.StringVar(&cfg.Addr, "a", cfg.Addr, "адрес эндпоинта HTTP-сервера")

    err := flagSet.Parse(os.Args[1:])
    if err != nil {
        return nil, err
    }

    // Проверяем, откуда взято значение адреса
    if cfg.Addr == "localhost:8080" {
        log.Println("Переменная окружения ADDRESS не найдена, используется значение по умолчанию:", cfg.Addr)
    } else if flagSet.Lookup("a") != nil && flagSet.Lookup("a").Value.String() != "" {
        log.Println("Используется флаг командной строки -a:", cfg.Addr)
    } else {
        log.Println("Используется переменная окружения ADDRESS:", cfg.Addr)
    }

    return cfg, nil
}