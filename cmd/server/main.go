// Объявляем пакет main для текущего файла
package main

// Импортируем необходимые пакеты
import (
    "log"     // Для логирования событий
    "net/http" // Для работы с HTTP-запросами
    "strconv"  // Для преобразования строк в числа
    "strings"  // Для работы со строками
    "sync"     // Для синхронизации доступа к общим ресурсам
)

// MemStorage структура для хранения метрик в памяти
type MemStorage struct {
    // Мьютекс для синхронизации доступа к данным
    mu sync.RWMutex
    // Карта для хранения gauge-метрик
    gauges map[string]float64
    // Карта для хранения counter-метрик
    counters map[string]int64
}

// Storage интерфейс для взаимодействия с хранилищем метрик
type Storage interface {
    // Метод для обновления gauge-метрики
    UpdateGauge(name string, value float64)
    // Метод для обновления counter-метрики
    UpdateCounter(name string, value int64)
}

// NewMemStorage конструктор для создания нового хранилища
func NewMemStorage() *MemStorage {
    // Возвращаем новое хранилище с инициализированными картами
    return &MemStorage{
        gauges:   make(map[string]float64), // Инициализируем карту для gauge-метрик
        counters: make(map[string]int64),  // Инициализируем карту для counter-метрик
    }
}

// UpdateGauge обновляет значение gauge-метрики
func (s *MemStorage) UpdateGauge(name string, value float64) {
    // Блокируем мьютекс для записи
    s.mu.Lock()
    // Отменяем блокировку при выходе из функции
    defer s.mu.Unlock()
    // Обновляем значение метрики
    s.gauges[name] = value
}

// UpdateCounter обновляет значение counter-метрики
func (s *MemStorage) UpdateCounter(name string, value int64) {
    // Блокируем мьютекс для записи
    s.mu.Lock()
    // Отменяем блокировку при выходе из функции
    defer s.mu.Unlock()
    // Добавляем значение к существующему или создаем новую метрику
    s.counters[name] += value
}

// UpdateHandler обработчик запросов на обновление метрик
func UpdateHandler(storage Storage) http.HandlerFunc {
    // Возвращаем функцию-обработчик
    return func(w http.ResponseWriter, r *http.Request) {
        // Устанавливаем заголовок Content-Type для ответа
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")

        // Проверяем, что запрос был отправлен методом POST
        if r.Method != http.MethodPost {
            // Если метод не POST, возвращаем статус 405
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        // Разбиваем URL на части
        pathParts := strings.Split(r.URL.Path, "/")[2:]
        // Проверяем, что URL имеет правильный формат
        if len(pathParts) != 3 {
            // Если формат неверный, возвращаем статус 404
            w.WriteHeader(http.StatusNotFound)
            return
        }

        // Извлекаем тип метрики из URL
        metricType := pathParts[0]
        // Извлекаем имя метрики из URL
        metricName := pathParts[1]
        // Извлекаем значение метрики из URL
        metricValue := pathParts[2]

        // Проверяем, что имя метрики не пустое
        if metricName == "" {
            // Если имя пустое, возвращаем статус 404
            w.WriteHeader(http.StatusNotFound)
            return
        }

        // Обрабатываем метрики в зависимости от типа
        switch metricType {
        case "gauge":
            // Попытка парсить значение как float64
            value, err := strconv.ParseFloat(metricValue, 64)
            // Если парсинг не удался, возвращаем статус 400
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            // Обновляем значение gauge-метрики
            storage.UpdateGauge(metricName, value)

        case "counter":
            // Попытка парсить значение как int64
            value, err := strconv.ParseInt(metricValue, 10, 64)
            // Если парсинг не удался, возвращаем статус 400
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            // Обновляем значение counter-метрики
            storage.UpdateCounter(metricName, value)

        default:
            // Если тип метрики не поддерживается, возвращаем статус 400
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // Если все прошло успешно, возвращаем статус 200
        w.WriteHeader(http.StatusOK)
    }
}

// Основная функция программы
func main() {
    // Создаем новое хранилище метрик
    storage := NewMemStorage()
    // Создаем новый мультиплексор для маршрутизации запросов
    mux := http.NewServeMux()
    // Регистрируем обработчик для маршрута "/update/"
    mux.Handle("/update/", UpdateHandler(storage))

    // Создаем новый сервер
    server := &http.Server{
        // Устанавливаем адрес сервера
        Addr:    "localhost:8080",
        // Устанавливаем обработчик запросов
        Handler: mux,
    }

    // Логируем сообщение о запуске сервера
    log.Printf("Starting server on http://localhost:8080\n")
    // Запускаем сервер и ждем его завершения
    log.Fatal(server.ListenAndServe())
}