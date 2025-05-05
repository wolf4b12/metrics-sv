package agentmethods

import (
    "bytes"
    "compress/gzip"
    "fmt"
)

// compressPayload сжимает данные с помощью Gzip
func (a *Agent) CompressPayload(data []byte) ([]byte, error) {
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
