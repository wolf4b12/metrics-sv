package agentmethods

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

// sendSingleTextMetric отправляет одну метрику в текстовом формате
func (a *Agent) sendSingleTextMetric(
	urlPath string,
	payload string,
	metricID string,
	checkRequired func() bool,
) {
	if !checkRequired() {
		log.Printf("Отсутствует обязательное поле для метрики '%s'\n", metricID)
		return
	}

	// Сжимаем payload
	compressedData, err := a.compressPayload([]byte(payload))
	if err != nil {
		a.handleErrorAndContinue("сжатия URL", err)
		return
	}

	// Формируем POST-запрос с Gzip-данными
	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewBuffer(compressedData))
	if err != nil {
		a.handleErrorAndContinue("формирования запроса", err)
		return
	}

	// Устанавливаем заголовки
	a.SetHeaders(req, "text/plain")

	// Выполняем запрос
	resp, err := a.client.Do(req)
	if err != nil {
		a.handleErrorAndContinue("отправки метрики", err)
		return
	}

	// Обрабатываем ответ
	if err := a.handleResponse(resp); err != nil {
		a.handleErrorAndContinue("обработки ответа", err)
	}
}

// SendTextCollectedMetrics отправляет собранные метрики в текстовом формате
func (a *Agent) SendTextCollectedMetrics() {
	for {
		a.mu.Lock()

		baseURL := fmt.Sprintf("http://%s/update", a.addr)

		for _, gauge := range a.Gauges {
			a.sendSingleTextMetric(
				baseURL+"/gauge",
				fmt.Sprintf("%s/gauge/%s/%f", baseURL, gauge.ID, *(gauge.Value)),
				gauge.ID,
				func() bool { return gauge.Value != nil },
			)
		}

		for _, counter := range a.Counters {
			a.sendSingleTextMetric(
				baseURL+"/counter",
				fmt.Sprintf("%s/counter/%s/%d", baseURL, counter.ID, *(counter.Delta)),
				counter.ID,
				func() bool { return counter.Delta != nil },
			)
		}

		a.mu.Unlock()
		time.Sleep(a.reportInterval)
	}
}
