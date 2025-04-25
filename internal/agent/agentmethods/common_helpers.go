package agentmethods

import (
    "bytes"
    "compress/gzip"
    "fmt"
    "io"
    "log"
    "net/http"
)

// compressPayload сжимает данные с помощью Gzip
func (a *Agent) compressPayload(data []byte) ([]byte, error) {
    var buf bytes.Buffer
    zw := gzip.NewWriter(&buf)
    if _, err := zw.Write(data); err != nil {
        return nil, fmt.Errorf("ошибка сжатия метрики: %v", err)
    }
    if err := zw.Close(); err != nil {
        return nil, fmt.Errorf("ошибка закрытия компрессора: %v", err)
    }
    return buf.Bytes(), nil
}

// handleResponse обрабатывает ответ от сервера
func (a *Agent) handleResponse(resp *http.Response) error {
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
func (a *Agent) handleErrorAndContinue(action string, err error) {
    log.Printf("Ошибка %s: %v\n", action, err)
}