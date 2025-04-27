package config

import (
    "flag"
    "log"
    "os"
    "time"
    "reflect"

    "github.com/caarlos0/env/v11"
)

type Config struct {
    Addr              string        `env:"ADDRESS" envDefault:"localhost:8080"`                 // Адрес слушания
    StoreIntervalSec  int           `env:"STORE_INTERVAL" envDefault:"300"`                      // Интервал автосохранения в секундах
    RestoreOnStartup  bool          `env:"RESTORE" envDefault:"false"`                         // Восстанавливать метрики при старте
    FileStoragePath   string        `env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics.json"`   // Путь к файлу хранения метрик
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

// Создаем новый экземпляр конфигурации
func NewConfig() (*Config, error) {
    cfg := &Config{}

    // Читаем переменные окружения
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }

    // Сохраняем первоначальные значения по умолчанию
    defaultValues := map[string]interface{}{
        "a": cfg.Addr,
        "i": cfg.StoreIntervalSec,
        "f": cfg.FileStoragePath,
        "r": cfg.RestoreOnStartup,
    }

    // Формируем набор флагов командной строки
    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    flagSet.StringVar(&cfg.Addr, "a", cfg.Addr, "Адрес HTTP-эндпоинта")
    flagSet.IntVar(&cfg.StoreIntervalSec, "i", cfg.StoreIntervalSec, "Интервал автосохранения в секундах")
    flagSet.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "Путь к файлу хранения метрик")
    flagSet.BoolVar(&cfg.RestoreOnStartup, "r", cfg.RestoreOnStartup, "Включение восстановления метрик при старте")

    err := flagSet.Parse(os.Args[1:])
    if err != nil {
        return nil, err
    }

    // Проверка источника для поля Addr
    checkSource(flagSet, defaultValues["a"], "a", "ADDRESS", cfg.Addr)

    // Остальные поля проверяются аналогично
    checkSourceInt(flagSet, defaultValues["i"].(int), "i", "STORE_INTERVAL", cfg.StoreIntervalSec)
    checkSourceString(flagSet, defaultValues["f"].(string), "f", "FILE_STORAGE_PATH", cfg.FileStoragePath)
    checkSourceBool(flagSet, defaultValues["r"].(bool), "r", "RESTORE", cfg.RestoreOnStartup)

    return cfg, nil
}

// Вспомогательная функция для вывода источника установки значения
func checkSource(flagSet *flag.FlagSet, defValue interface{}, flagName, envKey string, currentValue interface{}) {
    envVal := os.Getenv(envKey)
    if envVal != "" { // Установлена переменная окружения
        log.Printf("Используется переменная окружения %s: %+v\n", envKey, currentValue)
    } else if flagSet.Lookup(flagName) != nil && !reflect.DeepEqual(currentValue, defValue) { // Явно изменённый флаг
        log.Printf("Используется флаг командной строки -%s: %+v\n", flagName, currentValue)
    } else { // Ничего не передано, использовано значение по умолчанию
        log.Printf("%s не установлено, используется значение по умолчанию: %+v\n", envKey, currentValue)
    }
}

// Специфическая реализация для Int типов
func checkSourceInt(flagSet *flag.FlagSet, defValue int, flagName, envKey string, currentValue int) {
    envVal := os.Getenv(envKey)
    if envVal != "" { // Установлена переменная окружения
        log.Printf("Используется переменная окружения %s: %+v\n", envKey, currentValue)
    } else if flagSet.Lookup(flagName) != nil && currentValue != defValue { // Явно изменённый флаг
        log.Printf("Используется флаг командной строки -%s: %+v\n", flagName, currentValue)
    } else { // Ничего не передано, использовано значение по умолчанию
        log.Printf("%s не установлено, используется значение по умолчанию: %+v\n", envKey, currentValue)
    }
}

// Специфическая реализация для String типов
func checkSourceString(flagSet *flag.FlagSet, defValue string, flagName, envKey string, currentValue string) {
    envVal := os.Getenv(envKey)
    if envVal != "" { // Установлена переменная окружения
        log.Printf("Используется переменная окружения %s: %+v\n", envKey, currentValue)
    } else if flagSet.Lookup(flagName) != nil && currentValue != defValue { // Явно изменённый флаг
        log.Printf("Используется флаг командной строки -%s: %+v\n", flagName, currentValue)
    } else { // Ничего не передано, использовано значение по умолчанию
        log.Printf("%s не установлено, используется значение по умолчанию: %+v\n", envKey, currentValue)
    }
}

// Реализация для Bool типа
func checkSourceBool(flagSet *flag.FlagSet, defValue bool, flagName, envKey string, currentValue bool) {
    strValue := "false"
    if currentValue {
        strValue = "true"
    }

    envVal := os.Getenv(envKey)
    if envVal != "" { // Установлена переменная окружения
        log.Printf("Используется переменная окружения %s: %s\n", envKey, strValue)
    } else if flagSet.Lookup(flagName) != nil && currentValue != defValue { // Явно изменённый флаг
        log.Printf("Используется флаг командной строки -%s: %s\n", flagName, strValue)
    } else { // Ничего не передано, использовано значение по умолчанию
        log.Printf("%s не установлено, используется значение по умолчанию: %s\n", envKey, strValue)
    }
}