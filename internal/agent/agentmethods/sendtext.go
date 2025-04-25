package agentmethods

import (
    "fmt"
    "log"
    "time"


)


// SendTextCollectedMetrics отправляет собранные метрики в текстовом формате
func (a *Agent) SendTextCollectedMetrics() {
    for {
        a.mu.Lock()

        // Подготовим URL для отправки метрик
        baseURL := fmt.Sprintf("http://%s/update", a.addr)

        // Отправляем измерители (Gauges)
        for _, gauge := range a.Gauges {
            if gauge.Value == nil {
                log.Printf("Отсутствует обязательное поле 'Value' для датчика '%s'\n", gauge.ID)
                continue
            }

            // Формируем URL для конкретной метрики
            textURL := fmt.Sprintf("%s/gauge/%s/%f", baseURL, gauge.ID, *(gauge.Value))

            // Отправляем метрику
            if err := a.sendMetric(baseURL+"/gauge", []byte(textURL), "text/plain"); err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
            }
        }

        // Отправляем счетчики (Counters)
        for _, counter := range a.Counters {
            if counter.Delta == nil {
                log.Printf("Отсутствует обязательное поле 'Delta' для счётчика '%s'\n", counter.ID)
                continue
            }

            // Формируем URL для конкретной метрики
            textURL := fmt.Sprintf("%s/counter/%s/%d", baseURL, counter.ID, *(counter.Delta))

            // Отправляем метрику
            if err := a.sendMetric(baseURL+"/counter", []byte(textURL), "text/plain"); err != nil {
                log.Printf("Ошибка отправки метрики: %v\n", err)
            }
        }

        a.mu.Unlock()

        // Ждем указанный интервал
        time.Sleep(a.reportInterval)
    }
}