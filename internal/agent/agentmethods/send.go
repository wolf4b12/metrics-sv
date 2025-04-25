package agentmethods

import (
    "bytes"
    "fmt"
    "net/http"
    "compress/gzip"
    "io"

)


// Метод для отправки собранных метрик
// Метод для отправки собранных метрик в формате JSON с поддержкой Gzip

// sendMetric отправляет единичную метрику на сервер
func (a *Agent) sendMetric(url string, payload []byte, contentType string) error {
    // Сжимаем данные с помощью Gzip
    var buf bytes.Buffer
    zw := gzip.NewWriter(&buf)
    if _, err := zw.Write(payload); err != nil {
        return fmt.Errorf("ошибка сжатия метрики: %v", err)
    }
    if err := zw.Close(); err != nil {
        return fmt.Errorf("ошибка закрытия компрессора: %v", err)
    }

    // Формируем запрос с Gzip-данными
    req, err := http.NewRequest(http.MethodPost, url, &buf)
    if err != nil {
        return fmt.Errorf("ошибка формирования запроса: %v", err)
    }

    // Устанавливаем заголовки для сжатия
    req.Header.Set("Content-Type", contentType)
    req.Header.Set("Content-Encoding", "gzip")
    req.Header.Set("Accept-Encoding", "gzip")

    // Выполняем запрос
    resp, err := a.client.Do(req)
    if err != nil {
        return fmt.Errorf("ошибка отправки метрики: %v", err)
    }
    defer resp.Body.Close()

    // Проверяем статус ответа
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("получен неправильный статус-код (%d)", resp.StatusCode)
    }

    // Если ответ приходит в сжатом виде, разархивируем его
    if resp.Header.Get("Content-Encoding") == "gzip" {
        reader, err := gzip.NewReader(resp.Body)
        if err != nil {
            return fmt.Errorf("ошибка разбора Gzip-ответа: %v", err)
        }
        defer reader.Close()

        // Читаем ответ
        bodyBytes, err := io.ReadAll(reader)
        if err != nil {
            return fmt.Errorf("ошибка чтения тела ответа: %v", err)
        }

        fmt.Println(string(bodyBytes))
    } else {
        // Ответ несжатый, читаем обычный
        bodyBytes, err := io.ReadAll(resp.Body)
        if err != nil {
            return fmt.Errorf("ошибка чтения тела ответа: %v\n", err)
        }

        fmt.Println(string(bodyBytes))
    }

    return nil
}