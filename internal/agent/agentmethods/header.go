package agentmethods

import (
	"net/http"

)

// setHeaders устанавливает стандартные заголовки для запросов
func (a *Agent) SetHeaders(req *http.Request, contentType string) {
    req.Header.Set("Content-Type", contentType)
    req.Header.Set("Content-Encoding", "gzip")
    req.Header.Set("Accept-Encoding", "gzip")
}