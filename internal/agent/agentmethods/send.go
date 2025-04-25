package agentmethods

import (
    "fmt"
    "net/http"
    "compress/gzip"
    "io"

)


// doRequest отправляет HTTP-запрос и обрабатывает ответ
func (a *Agent) doRequest(method, url string, body io.Reader, headers map[string]string) error {
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return fmt.Errorf("ошибка формирования запроса: %v", err)
    }

    // Установка заголовков
    for k, v := range headers {
        req.Header.Set(k, v)
    }

    // Выполнение запроса
    resp, err := a.client.Do(req)
    if err != nil {
        return fmt.Errorf("ошибка отправки метрики: %v", err)
    }
    defer resp.Body.Close()

    // Проверка статуса ответа
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("получен неправильный статус-код (%d)", resp.StatusCode)
    }

    // Обработка ответа
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
            return fmt.Errorf("ошибка чтения тела ответа: %v", err)
        }

        fmt.Println(string(bodyBytes))
    }

    return nil
}