package pg

import (
    "database/sql"
    "encoding/json"
    "log"
    "strings"
    "sync"
    "github.com/wolf4b12/metrics-sv/internal/server/storage"

    _ "github.com/lib/pq" // Импорт драйвера для PostgreSQL
    
)

// NewPGStorage — strore для работы с PostgreSQL
type PGStore struct {
    db     *sql.DB
    mu     sync.Mutex
}

// NewPGStorage создает новый store для работы с PostgreSQL
func NewPGStorage(connStr string) (*PGStore, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    pga := &PGStore{
        db: db,
    }

    // Автоматически создаем таблицу, если она отсутствует
    err = pga.createTables()
    if err != nil {
        return nil, err
    }

    return pga, nil
}

// createTables проверяет наличие таблиц и создает их при отсутствии
func (pga *PGStore) createTables() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS gauges (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) UNIQUE NOT NULL,
            value DOUBLE PRECISION NOT NULL
        );`,
        `CREATE TABLE IF NOT EXISTS counters (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) UNIQUE NOT NULL,
            value BIGINT NOT NULL
        );`,
    }

    for _, q := range queries {
        _, err := pga.db.Exec(q)
        if err != nil {
            return err
        }
    }

    return nil
}

// Set устанавливает значение по ключу
func (pga *PGStore) Set(key string, value any) {
   _, err := json.Marshal(value)
    if err != nil {
        log.Printf("Ошибка при маршалинге значения: %v\n", err)
        return
    }

    tx, err := pga.db.Begin()
    if err != nil {
        log.Printf("Ошибка начала транзакции: %v\n", err)
        return
    }
    defer tx.Rollback()

    switch val := value.(type) {
    case storage.MetricValue:
        switch val.Type {
        case storage.Gauge:
            _, err = tx.Exec(`
                INSERT INTO gauges (name, value) VALUES($1, $2)
                ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;
            `, key, val.Value)
        case storage.Counter:
            _, err = tx.Exec(`
                INSERT INTO counters (name, value) VALUES($1, $2)
                ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;
            `, key, val.Value)
        }
        if err != nil {
            log.Printf("Ошибка при сохранении значения: %v\n", err)
            return
        }
    }

    err = tx.Commit()
    if err != nil {
        log.Printf("Ошибка при завершении транзакции: %v\n", err)
    }
}

// Get получает значение по ключу
func (pga *PGStore) Get(key string) (any, bool) {
    type result struct {
        Name  string
        Value any
    }

    stmt := `
    SELECT 'gauge' AS type, g.name, g.value::FLOAT FROM gauges g WHERE g.name = $1
    UNION ALL
    SELECT 'counter', c.name, c.value::BIGINT FROM counters c WHERE c.name = $1
    LIMIT 1;
    `

    var res result
    row := pga.db.QueryRow(stmt, key)
    err := row.Scan(&res.Name, &res.Value)
    if err != nil && err != sql.ErrNoRows {
        log.Printf("Ошибка при извлечении значения: %v\n", err)
        return nil, false
    }

    if strings.TrimSpace(res.Name) == "" { // Если ничего не нашли
        return nil, false
    }

    return storage.MetricValue{Type: storage.MetricType(res.Name), Value: res.Value}, true
}

// Delete удаляет значение по ключу
func (pga *PGStore) Delete(key string) {
    tx, err := pga.db.Begin()
    if err != nil {
        log.Printf("Ошибка начала транзакции: %v\n", err)
        return
    }
    defer tx.Rollback()

    _, err = tx.Exec("DELETE FROM gauges WHERE name = $1;", key)
    if err != nil {
        log.Printf("Ошибка при удалении gauge-метрики: %v\n", err)
        return
    }

    _, err = tx.Exec("DELETE FROM counters WHERE name = $1;", key)
    if err != nil {
        log.Printf("Ошибка при удалении counter-метрики: %v\n", err)
        return
    }

    err = tx.Commit()
    if err != nil {
        log.Printf("Ошибка при завершении транзакции: %v\n", err)
    }
}

// All возвращает все данные
func (pga *PGStore) All() map[string]any {
    results := make(map[string]any)

    gQuery := "SELECT name, value FROM gauges;"
    cQuery := "SELECT name, value FROM counters;"

    rows, err := pga.db.Query(gQuery)
    if err != nil {
        log.Printf("Ошибка при извлечении gauge-метрик: %v\n", err)
        return results
    }
    defer rows.Close()

    for rows.Next() {
        var name string
        var value float64
        err := rows.Scan(&name, &value)
        if err != nil {
            log.Printf("Ошибка при обработке результата: %v\n", err)
            continue
        }
        results[name] = storage.MetricValue{Type: storage.Gauge, Value: value}
    }

    // Обязательно проверяем наличие ошибок после цикла Next()
    if err := rows.Err(); err != nil {
        log.Printf("Ошибка при обработке gauge-метрик: %v\n", err)
    }

    rows, err = pga.db.Query(cQuery)
    if err != nil {
        log.Printf("Ошибка при извлечении counter-метрик: %v\n", err)
        return results
    }
    defer rows.Close()

    for rows.Next() {
        var name string
        var value int64
        err := rows.Scan(&name, &value)
        if err != nil {
            log.Printf("Ошибка при обработке результата: %v\n", err)
            continue
        }
        results[name] = storage.MetricValue{Type: storage.Counter, Value: value}
    }

    // Опять проверяем наличие ошибок после цикла Next()
    if err := rows.Err(); err != nil {
        log.Printf("Ошибка при обработке counter-метрик: %v\n", err)
    }

    return results
}