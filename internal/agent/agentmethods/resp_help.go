package agentmethods

import (
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "net/http"
)


// handleResponse обрабатывает ответ от сервера
func (a *Agent) HandleResponse(resp *http.Response) error {
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
            return fmt.Errorf("ошибка чтения тела ответа: %v", err)
        }

        fmt.Println(string(bodyBytes))
    }

    return nil
}

// handleErrorAndContinue обрабатывает ошибку и продолжает выполнение
func (a *Agent) HandleErrorAndContinue(action string, err error) {
    log.Printf("Ошибка %s: %v\n", action, err)
}