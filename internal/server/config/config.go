package config

import (
    "flag"
    "log"
    "os"
)


type Config struct {
    addr string
}


// Метод для получения адреса
func (c *Config) GetAddr() string {
    return c.addr
}


func NewConfig() *Config {
    cfg := &Config{}

    if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
        cfg.addr = envAddr
        log.Println("Используется переменная окружения ADDRESS:", cfg.addr)
        return cfg
    }

    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    defaultAddr := "localhost:8080"
    flagSet.StringVar(&cfg.addr, "a", defaultAddr, "адрес эндпоинта HTTP-сервера")

    err := flagSet.Parse(os.Args[1:])
    if err != nil {
        log.Fatalf("Ошибка парсинга флагов: %v", err)
    }

    if flagSet.NArg() > 0 {
        log.Fatalf("Неизвестный флаг: %s\n", flagSet.Arg(flagSet.NArg()-1))
    }

    if cfg.addr == "" {
        cfg.addr = defaultAddr
        log.Println("Переменная окружения ADDRESS не найдена, используется значение по умолчанию:", cfg.addr)
    } else {
        log.Println("Используется флаг командной строки -a:", cfg.addr)
    }

    return cfg
}