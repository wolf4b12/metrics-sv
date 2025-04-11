package config // Объявляем пакет 'config'.

import ( // Начинаем секцию импорта пакетов.
    "flag"     // Импортируем пакет для работы с флагами командной строки.
    "log"      // Импортируем пакет для ведения логирования.
    "os"       // Импортируем пакет для взаимодействия с операционной системой.
    "github.com/caarlos0/env/v11" // Импортируем библиотеку 'env' для парсинга переменных окружения.
)

type Config struct { // Определяем структуру 'Config' для хранения настроек.
    Addr string `env:"ADDRESS" envDefault:"localhost:8080"` // Поле 'Addr' хранит адрес, аннотировано для автоматического заполнения из переменной окружения 'ADDRESS' или флага '-a'.
}

// Метод для получения адреса
func (c *Config) GetAddr() string { // Метод возвращает значение поля 'Addr'.
    return c.Addr // Возвращаем текущее значение адреса.
}

func NewConfig() (*Config, error) { // Функция для создания новой конфигурации.
    cfg := &Config{} // Создаём новую структуру 'Config'.

    // Парсим переменные окружения
    if err := env.Parse(cfg); err != nil { // Пробуем распарсить переменные окружения.
        return nil, err // Если возникла ошибка, возвращаем nil и ошибку.
    }

    // Парсим флаги командной строки
    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Создаём новый набор флагов с именем программы.
    flagSet.StringVar(&cfg.Addr, "a", "", "адрес эндпоинта HTTP-сервера") // Добавляем флаг '-a' для задания адреса.

    err := flagSet.Parse(os.Args[1:]) // Парсим аргументы командной строки.
    if err != nil {
        return nil, err // Если возникла ошибка, возвращаем nil и ошибку.
    }

    // Проверяем, какой источник использовал адрес
    if cfg.Addr == "" || cfg.Addr == "localhost:8080" { // Если адрес пустой или равен значению по умолчанию...
        log.Println("Переменная окружения ADDRESS не найдена, используется значение по умолчанию:", cfg.Addr) // Логируем использование значения по умолчанию.
    } else if flagSet.Lookup("a") != nil && flagSet.Lookup("a").Value.String() != "" { // Если флаг '-a' был указан и имеет значение...
        log.Println("Используется флаг командной строки -a:", cfg.Addr) // Логируем использование флага командной строки.
    } else {
        log.Println("Используется переменная окружения ADDRESS:", cfg.Addr) // Логируем использование переменной окружения.
    }

    return cfg, nil // Возвращаем сконфигурированную структуру и отсутствие ошибок.
}