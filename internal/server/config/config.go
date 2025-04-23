package config

import (
    "flag"
    "log"
    "os"
//    "strconv"
//    "strings"
    "time"

    "github.com/caarlos0/env/v11"
)

type Config struct {
    Addr              string        `env:"ADDRESS" envDefault:"localhost:8080"`               // Адрес слушания
    StoreIntervalSec  int           `env:"STORE_INTERVAL" envDefault:"3"`                  // Интервал автосохранения в секундах
    RestoreOnStartup  bool          `env:"RESTORE" envDefault:"false"`                       // Восстанавливать метрики при старте
    FileStoragePath   string        `env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics.json"`     // Путь к файлу хранения метрик
}

// Метод для получения адреса
func (c *Config) GetAddr() string {
    return c.Addr
}

// Возвращает интервал автосохранения в виде Duration
func (c *Config) GetStoreInterval() time.Duration {
    return time.Duration(c.StoreIntervalSec) * time.Second
}

// Метод для удобного получения пути к файлу хранения
func (c *Config) GetFileStoragePath() string {
    return c.FileStoragePath
}

// Метод для проверки параметра restore-on-startup
func (c *Config) IsRestoreEnabled() bool {
    return c.RestoreOnStartup
}

// Создает новый экземпляр конфигурации
func NewConfig() (*Config, error) {
    cfg := &Config{}

    // Читаем переменные окружения
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }

    // Формируем набор флагов командной строки
    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    flagSet.StringVar(&cfg.Addr, "a", cfg.Addr, "адрес HTTP-эндпоинта")
    flagSet.IntVar(&cfg.StoreIntervalSec, "i", cfg.StoreIntervalSec, "интервал автосохранения в секундах")
    flagSet.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "путь к файлу хранения метрик")
    flagSet.BoolVar(&cfg.RestoreOnStartup, "r", cfg.RestoreOnStartup, "включить восстановление метрик при старте")

    err := flagSet.Parse(os.Args[1:])
    if err != nil {
        return nil, err
    }

    // Проверка, откуда взяты параметры
    if cfg.Addr == "localhost:8080" {
        log.Println("Переменная окружения ADDRESS не найдена, используется значение по умолчанию:", cfg.Addr)
    } else if flagSet.Lookup("a") != nil && flagSet.Lookup("a").Value.String() != "" {
        log.Println("Используется флаг командной строки -a:", cfg.Addr)
    } else {
        log.Println("Используется переменная окружения ADDRESS:", cfg.Addr)
    }

    // Аналогично для остальных полей
    checkSource(flagSet, "i", "STORE_INTERVAL", cfg.StoreIntervalSec)
    checkSource(flagSet, "f", "FILE_STORAGE_PATH", cfg.FileStoragePath)
    checkSourceBool(flagSet, "r", "RESTORE", cfg.RestoreOnStartup)

    return cfg, nil
}

// Вспомогательная функция для вывода источника установки значения
func checkSource(flagSet *flag.FlagSet, flagName, envKey string, value interface{}) {
    if flagSet.Lookup(flagName) != nil && flagSet.Lookup(flagName).Value.String() != "" {
        log.Printf("Используется флаг командной строки -%s: %+v\n", flagName, value)
    } else if os.Getenv(envKey) != "" {
        log.Printf("Используется переменная окружения %s: %+v\n", envKey, value)
    } else {
        log.Printf("%s не установлено, используется значение по умолчанию: %+v\n", envKey, value)
    }
}

// Аналогичная проверка для булевых параметров
func checkSourceBool(flagSet *flag.FlagSet, flagName, envKey string, value bool) {
    strValue := "false"
    if value {
        strValue = "true"
    }

    if flagSet.Lookup(flagName) != nil && flagSet.Lookup(flagName).Value.String() != "" {
        log.Printf("Используется флаг командной строки -%s: %s\n", flagName, strValue)
    } else if os.Getenv(envKey) != "" {
        log.Printf("Используется переменная окружения %s: %s\n", envKey, strValue)
    } else {
        log.Printf("%s не установлено, используется значение по умолчанию: %s\n", envKey, strValue)
    }
}