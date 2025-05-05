package main

import (
    "math/rand"
    "time"
    "sync"
    "github.com/wolf4b12/metrics-sv/internal/agent/agentmethods" // Импортируем пакет agentmethods
    "github.com/wolf4b12/metrics-sv/internal/agent/parseflags" // Импортируем пакет parseflags
)

func main() {
    rand.New(rand.NewSource(time.Now().UnixNano())) // Создаем источник случайных чисел

    poll, report, addr := parseflags.ParseFlags()

    agent := agentmethods.NewAgent(poll, report, addr)

    // Глобальный mutex для защиты общей структуры
    var globalMu sync.Mutex

    // Функции-колбеки для безопасной работы с общим состоянием
    collectMetrics := func() {
        globalMu.Lock()
        defer globalMu.Unlock()
        agent.StartCollectingMetrics()
    }
    sendJSONMetrics := func() {
        globalMu.Lock()
        defer globalMu.Unlock()
        agent.SendJSONCollectedMetrics()
    }
    sendTextMetrics := func() {
        globalMu.Lock()
        defer globalMu.Unlock()
        agent.SendTextCollectedMetrics()
    }
    sendBatches := func() {
        globalMu.Lock()
        defer globalMu.Unlock()
        agent.CollectAndSendBatches()
    }

    // Горутины с безопасным доступом к общему ресурсу
    go collectMetrics()
    go sendJSONMetrics()
    go sendTextMetrics()
    go sendBatches()

    // Основной поток ждет бесконечно
    select {}
}